package sigs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/benpate/derp"
	domaintools "github.com/benpate/domain"
	"github.com/benpate/hannibal"
	"github.com/benpate/rosetta/slice"
	"github.com/rs/zerolog/log"
)

// Signer contains all of the settings necessary to sign a request
type Signer struct {
	PublicKeyID   string
	PrivateKey    crypto.PrivateKey
	Fields        []string
	SignatureHash crypto.Hash
	BodyDigest    crypto.Hash
	HS2019        bool
	Created       int64
	Expires       int64
}

// NewSigner returns a fully initialized Signer
func NewSigner(publicKeyID string, privateKey crypto.PrivateKey, options ...SignerOption) Signer {
	result := Signer{
		PublicKeyID:   publicKeyID,
		PrivateKey:    privateKey,
		Fields:        []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		SignatureHash: crypto.SHA256,
		BodyDigest:    crypto.SHA256,
		Created:       0,
		Expires:       0,
	}
	result.With(options...)
	return result
}

// Sign signs a given http.Request.  It is syntactic
// sugar for NewSigner(options...).Sign(request)
func Sign(request *http.Request, publicKeyID string, privateKey crypto.PrivateKey, options ...SignerOption) error {
	signer := NewSigner(publicKeyID, privateKey, options...)
	return signer.Sign(request)
}

// Use applies the given options to the Signer
func (signer *Signer) With(options ...SignerOption) {
	for _, option := range options {
		option(signer)
	}
}

// Sign generates a signature and applies it to the given http.Request
func (signer *Signer) Sign(request *http.Request) error {

	// Try to generate a signature
	signature, err := signer.MakeSignature(request)

	if err != nil {
		return derp.Wrap(err, "hannibal.sigs.Sign", "Error getting signature")
	}

	signature.Created = signer.Created
	signature.Expires = signer.Expires

	// Apply the signature to the request
	request.Header.Set("Signature", signature.String())

	// Success. Can you imagine a more beautiful thing?
	return nil

}

// MakeSignature generates a Signature string for the given http.Request
func (signer *Signer) MakeSignature(request *http.Request) (Signature, error) {

	const location = "hannibal.sigs.MakeSignature"

	signature := NewSignature()

	// Only apply the digest if the "digest" field is in use
	if slice.Contains(signer.Fields, FieldDigest) {

		// Select a Digest Function
		digestName := getDigestName(signer.BodyDigest)
		digestFunc, err := getDigestFunc(signer.BodyDigest)

		if err != nil {
			return Signature{}, derp.Wrap(err, location, "Unable to create digest function")
		}

		// Apply the Digest function to the body
		if err := ApplyDigest(request, digestName, digestFunc); err != nil {
			return Signature{}, derp.Wrap(err, location, "Error applying digest")
		}
	}

	// If "date" field is in use, then verify that it's present in the header.
	// If the "date" field is invalid or unset, use the current time.
	if slice.Contains(signer.Fields, FieldDate) {
		date := request.Header.Get(FieldDate)
		if _, err := time.Parse(http.TimeFormat, date); err != nil {
			request.Header.Set(FieldDate, hannibal.TimeFormat(time.Now()))
		}
	}

	// Assemble the plaintext string from the configured request fields
	plainText := makePlaintext(request, signature, signer.Fields...)

	// Create a digest of the plaintext string using the configured digest algorithm
	digestText, err := makeSignatureHash(plainText, signer.SignatureHash)

	if err != nil {
		return Signature{}, derp.Wrap(err, location, "Unable to create digest")
	}

	// Sign the digest using the private key
	signedDigest, err := makeSignedDigest(digestText, signer.SignatureHash, signer.PrivateKey)

	if err != nil {
		return Signature{}, derp.Wrap(err, location, "Error signing digest")
	}

	// Assemble and return the signature object
	signature.KeyID = signer.PublicKeyID
	signature.Headers = signer.Fields
	signature.Algorithm = getAlgorithmName(signer.PrivateKey, signer.SignatureHash)
	signature.Signature = signedDigest

	return signature, nil
}

/******************************************
 * Helper Functions
 ******************************************/

// makePlaintext retrieves all fields from an HTTP request
func makePlaintext(request *http.Request, signature Signature, fields ...string) string {

	resultSlice := make([]string, len(fields))

	for index, field := range fields {
		value := getField(request, signature, field)
		resultSlice[index] = strings.ToLower(field) + ": " + value
	}

	// Join all fields together with a newline
	result := strings.Join(resultSlice, "\n")

	// Return the result (with logging)
	log.Trace().Str("plaintext", result).Msg("hannibal.sigs.makePlaintext")

	return result
}

// makeSignatureHash creates a digest of the provided plaintext string using the given digest algorithm
func makeSignatureHash(plaintext string, digestAlgorithm crypto.Hash) ([]byte, error) {

	const location = "hannibal.sigs.makeSignatureHash"

	var h hash.Hash
	var result []byte

	switch digestAlgorithm {

	case crypto.SHA256:
		h = sha256.New()

	case crypto.SHA512:
		h = sha512.New()

	default:
		return nil, derp.InternalError(location, "Unknown digest algorithm. Only sha-256 and sha-512 are supported", digestAlgorithm.String())
	}

	h.Write([]byte(plaintext))
	result = h.Sum(nil)

	if canTrace() {
		log.Trace().Str("loc", "hannibal.sigs.makeSignatureHash").Str("result", base64.StdEncoding.EncodeToString(result)).Send()
	}

	return result, nil
}

// makeSignedDigest signs the given digest using the provided private key.  It returns
// an error if the private key is not an RSA or ECDSA key.
func makeSignedDigest(digest []byte, hash crypto.Hash, privateKey crypto.PrivateKey) ([]byte, error) {

	const location = "hannibal.sigs.makeSignedDigest"

	switch typedValue := privateKey.(type) {

	case *rsa.PrivateKey:
		if resultBytes, err := rsa.SignPKCS1v15(rand.Reader, typedValue, hash, digest); err != nil {
			return nil, derp.Wrap(err, location, "Error signing hash with RSA private key")
		} else {
			return resultBytes, nil
		}

	case *ecdsa.PrivateKey:
		if resultBytes, err := ecdsa.SignASN1(rand.Reader, typedValue, digest); err != nil {
			return nil, derp.Wrap(err, location, "Error signing hash with ECDSA private key")
		} else {
			return resultBytes, nil
		}
	}

	return nil, derp.InternalError(location, "Unrecognized private key type", privateKey)
}

// getField retrieves the value of a named field from an HTTP request.
// It handles special cases for (request-target), (created), and (expires)
// fields, which are not stored in the HTTP header.
func getField(request *http.Request, signature Signature, field string) string {

	field = strings.Trim(field, " ")
	field = strings.ToLower(field)

	switch field {

	// Special case for (request-target) which needs to read the request body
	case FieldRequestTarget:
		return strings.ToLower(request.Method) + " " + getPathAndQuery(request.URL)

	// Special case for "host" which needs to read the request URL
	case FieldHost:
		return domaintools.TrueHostname(request)

	case FieldCreated:
		return signature.CreatedString()

	case FieldExpires:
		return signature.ExpiresString()
	}

	// All other fields are read from the http header
	// But let's do some extra work to make sure that multiple values are joined together
	// per: https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures#section-2.3
	fieldValue := request.Header[http.CanonicalHeaderKey(field)]
	return strings.Join(fieldValue, ", ")
}

// getPathAndQuery returns the path and query from a URL
func getPathAndQuery(url *url.URL) string {

	// NILCHECK
	if url == nil {
		return ""
	}

	result := url.Path

	if result == "" {
		result = "/"
	}

	if query := url.RawQuery; query != "" {
		result += "?" + query
	}

	return result
}

// getAlgorithmName returns the standard name used for the combination of private key and digest algorithms.
func getAlgorithmName(privateKey crypto.PrivateKey, digest crypto.Hash) string {

	var result string

	// Handle known private key types
	switch privateKey.(type) {

	case *rsa.PrivateKey:
		result = "rsa-"

	case *ecdsa.PrivateKey:
		result = "ecdsa-"

	default:
		// This is a fallback. It shouldn't happen.
		return Algorithm_HS2019
	}

	// Handle known digest hashes
	switch digest {

	case crypto.SHA256:
		result += "sha256"

	case crypto.SHA512:
		result += "sha512"

	default:
		// This is a fallback. It shouldn't happen.
		return Algorithm_HS2019
	}

	return result
}
