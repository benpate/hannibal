package sigs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
	"github.com/rs/zerolog/log"
)

// Signer contains all of the settings necessary to sign a request
type Signer struct {
	Fields        []string
	SignatureHash string
	BodyDigest    string
}

// NewSigner returns a fully initialized Signer
func NewSigner(options ...SignerOption) Signer {
	result := Signer{
		Fields:        []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		SignatureHash: Digest_SHA256,
		BodyDigest:    Digest_SHA256,
	}
	result.Use(options...)
	return result
}

// Sign signs a given http.Request.  It is syntactic
// sugar for NewSigner(options...).Sign(request)
func Sign(request *http.Request, privateKeyID string, privateKey crypto.PrivateKey, options ...SignerOption) error {
	signer := NewSigner(options...)
	return signer.Sign(request, privateKeyID, privateKey)
}

// Use applies the given options to the Signer
func (signer *Signer) Use(options ...SignerOption) {
	for _, option := range options {
		option(signer)
	}
}

// Sign signs the given http.Request
func (signer *Signer) Sign(request *http.Request, privateKeyID string, privateKey crypto.PrivateKey) error {

	// Add a body digest to the request
	digestFunc, err := getDigestFunc(signer.BodyDigest)

	if err != nil {
		return derp.Wrap(err, "hannibal.sigs.Sign", "Error creating digest function")
	}

	if err := ApplyDigest(request, digestFunc); err != nil {
		return derp.Wrap(err, "hannibal.sigs.Sign", "Error applying digest")
	}

	// If "date" field is in use, then verify that it's present in the header.
	// If the "date" field is invalid or unset, use the current time.
	if slice.Contains(signer.Fields, FieldDate) {
		date := request.Header.Get(FieldDate)
		if _, err := time.Parse(http.TimeFormat, date); err != nil {
			request.Header.Set(FieldDate, time.Now().Format(http.TimeFormat))
		}
	}

	// Assemble the plaintext string from the configured request fields
	plainText := makePlaintext(request, signer.Fields...)

	// Create a digest of the plaintext string using the configured digest algorithm
	digestText, err := makeSignatureHash(plainText, signer.SignatureHash)

	if err != nil {
		return derp.Wrap(err, "hannibal.sigs.Sign", "Error creating digest")
	}

	// Sign the digest using the private key
	signedDigest, err := makeSignedDigest(digestText, privateKey)

	if err != nil {
		return derp.Wrap(err, "hannibal.sigs.Sign", "Error signing digest")
	}

	// Assemble the signature object
	signature := NewSignature()
	signature.KeyID = privateKeyID
	signature.Headers = signer.Fields
	signature.Algorithm = Algorithm_HS2019
	signature.Signature = signedDigest

	// Add the signature to the request
	request.Header.Set("Signature", signature.String())

	// Success. Can you imagine a more beautiful thing?
	return nil
}

/******************************************
 * Helper Functions
 ******************************************/

// makePlaintext retrieves all fields from an HTTP request
func makePlaintext(request *http.Request, fields ...string) string {

	resultSlice := make([]string, len(fields))

	for index, field := range fields {
		value := getField(request, field)
		resultSlice[index] = strings.ToLower(field) + ": " + value
	}

	// Join all fields together with a newline
	result := strings.Join(resultSlice, "\n")

	// Return the result (with logging)
	log.Trace().Str("location", "hannibal.sigs.makePlaintext").Msg(result)
	return result
}

// makeSignatureHash creates a digest of the provided plaintext string using the given digest algorithm
func makeSignatureHash(plaintext string, digestAlgorithm string) ([]byte, error) {

	var result []byte

	switch digestAlgorithm {

	case Digest_SHA256:
		hash := sha256.New()
		hash.Write([]byte(plaintext))
		result = hash.Sum(nil)

	case Digest_SHA512:
		hash := sha512.New()
		hash.Write([]byte(plaintext))
		result = hash.Sum(nil)

	default:
		return nil, derp.New(500, "hannibal.sigs.hashPlaintext", "Unknown digest algorithm", digestAlgorithm)
	}

	log.Trace().Str("location", "hannibal.sigs.makeSignatureHash").Str("plaintext", plaintext).Str("result", string(result)).Send()
	return result, nil
}

// makeSignedDigest signs the given digest using the provided private key.  It returns
// an error if the private key is not an RSA or ECDSA key.
func makeSignedDigest(digest []byte, privateKey crypto.PrivateKey) ([]byte, error) {

	const location = "hannibal.sigs.makeSignedDigest"

	switch typedValue := privateKey.(type) {

	case *rsa.PrivateKey:
		if resultBytes, err := rsa.SignPKCS1v15(rand.Reader, typedValue, 0, digest); err != nil {
			return nil, derp.Wrap(err, location, "Error signing hash")
		} else {
			return resultBytes, nil
		}

	case *ecdsa.PrivateKey:
		if resultBytes, err := ecdsa.SignASN1(rand.Reader, typedValue, digest); err != nil {
			return nil, derp.Wrap(err, location, "Error signing hash")
		} else {
			return resultBytes, nil
		}
	}

	return nil, derp.NewInternalError(location, "Unrecognized private key type")
}

// getField retrieves the value of a named field from an HTTP request.
// It handles special cases for (request-target), (created), and (expires)
// fields, which are not stored in the HTTP header.
func getField(request *http.Request, field string) string {

	field = strings.Trim(field, " ")
	field = strings.ToLower(field)

	switch field {

	// Special case for (request-target) which needs to read the request body
	case FieldRequestTarget:
		return strings.ToLower(request.Method) + " " + getPathAndQuery(request.URL)

	// Special case for "host" which needs to read the request URL
	case FieldHost:
		return request.Host
	}

	// All other fields are read from the http header
	return request.Header.Get(field)
}

// getPathAndQuery returns the path and query from a URL
func getPathAndQuery(url *url.URL) string {

	result := url.Path

	if result == "" {
		result = "/"
	}

	if query := url.RawQuery; query != "" {
		result += "?" + query
	}

	return result
}
