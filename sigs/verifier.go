package sigs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
)

// Verifier contains all of the settings necessary to verify a request
type Verifier struct {
	Fields        []string
	SignatureHash string   // Digest algorithm used to create the signature.  Default is SHA256
	BodyDigests   []string // List of algorithms to accept from remote servers when they create a Digest header.  Default is SHA256 and SHA512
	Timeout       int      // Number of seconds before signatures are expired. Default is 43200 seconds (12 hours).
}

// NewVerifier returns a fully initialized Verifier
func NewVerifier(options ...VerifierOption) Verifier {
	result := Verifier{
		Fields:        []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		SignatureHash: Digest_SHA256,
		BodyDigests:   []string{Digest_SHA256, Digest_SHA512},
		Timeout:       12 * 60 * 60, // 12 hours
	}
	result.Use(options...)
	return result
}

// Verify verifies the given http.Request. This is
// syntactic sugar for NewVerifier(options...).Verify(request)
func Verify(request *http.Request, certificate string, options ...VerifierOption) error {
	verifier := NewVerifier(options...)
	return verifier.Verify(request, certificate)
}

// Use applies the given options to the Verifier
func (verifier *Verifier) Use(options ...VerifierOption) {
	for _, option := range options {
		option(verifier)
	}
}

// Verify verifies the given http.Request
func (verifier *Verifier) Verify(request *http.Request, certificate string) error {

	const location = "hannibal.sigs.Verify"

	// Verify the body Digest
	if err := VerifyDigest(request, verifier.BodyDigests...); err != nil {
		return derp.Wrap(err, location, "Error verifying body digest")
	}

	// Retrieve and parse the Signature from the HTTP Request
	signature, err := ParseSignature(request.Header.Get("Signature"))

	if err != nil {
		return derp.Wrap(err, location, "Error parsing signature")
	}

	// RULE: Verify that the signature contains all of the fields that we require
	if !slice.ContainsAll(signature.Headers, verifier.Fields...) {
		return derp.NewForbiddenError("hannibal.sigs.Verify", "Signature must include ALL of these fields", verifier.Fields)
	}

	// Recreate the plaintext and digest used to make the Signature
	plaintext := makePlaintext(request, signature.Headers...)
	digest, err := makeSignatureHash(plaintext, verifier.SignatureHash)

	if err != nil {
		return derp.Wrap(err, location, "Error creating digest")
	}

	// Decode the PEM certificate into a public key
	publicKey, err := DecodePublicPEM(certificate)

	if err != nil {
		return derp.Wrap(err, location, "Error decoding public key")
	}

	// Verify the signature matches the message digest
	if err := verifySignature(digest, signature.Signature, publicKey); err != nil {
		return derp.Wrap(err, location, "Invalid signature")
	}

	// Once I realzed it was successful, everything changed...
	return nil
}

/******************************************
 * Helper Functions
 ******************************************/

// verifySignature
func verifySignature(digest []byte, signature []byte, publicKey crypto.PublicKey) error {

	const location = "hannibal.sigs.verifySignature"

	switch typedValue := publicKey.(type) {

	case *rsa.PublicKey:
		if err := rsa.VerifyPKCS1v15(typedValue, 0, digest, signature); err != nil {
			return derp.Wrap(err, location, "Error verifying RSA signature")
		}
		return nil

	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(typedValue, digest, signature) {
			return derp.NewForbiddenError(location, "Invalid ECDSA signature")
		}
		return nil
	}

	return derp.NewBadRequestError("hannibal.sigs.verifySignature", "Unrecognized public key type")
}
