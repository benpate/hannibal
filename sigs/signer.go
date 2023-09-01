package sigs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/benpate/derp"
	"github.com/rs/zerolog/log"
)

// Signer contains all of the settings necessary to sign a request
type Signer struct {
	Fields          []string
	SignatureDigest string
	BodyDigests     []string
}

// NewSigner returns a fully initialized Signer
func NewSigner(options ...SignerOption) Signer {
	result := Signer{
		Fields:          []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		SignatureDigest: Digest_SHA256,
		BodyDigests:     []string{Digest_SHA256},
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

	// Assemble the plaintext string from the configured request fields
	plainText := makePlaintext(request, signer.Fields...)

	// Create a digest of the plaintext string using the configured digest algorithm
	digestText, err := makePlaintextDigest(plainText, signer.SignatureDigest)

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

// makePlaintextDigest creates a digest of the provided plaintext string using the given digest algorithm
func makePlaintextDigest(plaintext string, digestAlgorithm string) ([]byte, error) {

	switch digestAlgorithm {

	case Digest_SHA256:
		hash := sha256.New()
		hash.Write([]byte(plaintext))
		return hash.Sum(nil), nil

	case Digest_SHA512:
		hash := sha512.New()
		hash.Write([]byte(plaintext))
		return hash.Sum(nil), nil
	}

	return nil, derp.New(500, "hannibal.sigs.hashPlaintext", "Unknown digest algorithm", digestAlgorithm)
}

// makeSignedDigest signs the given digest using the provided private key.  It returns
// an error if the private key is not an RSA or ECDSA key.
func makeSignedDigest(digest []byte, privateKey crypto.PrivateKey) (string, error) {

	switch typedValue := privateKey.(type) {

	case *rsa.PrivateKey:
		if resultBytes, err := rsa.SignPKCS1v15(rand.Reader, typedValue, 0, digest); err != nil {
			return "", derp.Wrap(err, "hannibal.sigs.signHash", "Error signing hash")
		} else {
			return base64.StdEncoding.EncodeToString(resultBytes), nil
		}

	case *ecdsa.PrivateKey:
		if resultBytes, err := ecdsa.SignASN1(rand.Reader, typedValue, digest); err != nil {
			return "", derp.Wrap(err, "hannibal.sigs.signHash", "Error signing hash")
		} else {
			return base64.StdEncoding.EncodeToString(resultBytes), nil
		}
	}

	return "", derp.NewInternalError("hannibal.sigs.signHash", "Unrecognized private key type")
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
