package sigs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Verifier contains all of the settings necessary to verify a request
type Verifier struct {
	Fields          []string
	BodyDigests     []string // List of algorithms to accept from remote servers when they create a Digest header.  Default is SHA256 and SHA512
	SignatureHashes []string // Digest algorithm used to create the signature.  Default is SHA256, SHA512
	Timeout         int      // Number of seconds before signatures are expired. Default is 43200 seconds (12 hours).
}

// NewVerifier returns a fully initialized Verifier
func NewVerifier(options ...VerifierOption) Verifier {
	result := Verifier{
		Fields:          []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		BodyDigests:     []string{Digest_SHA256, Digest_SHA512},
		SignatureHashes: []string{Digest_SHA256, Digest_SHA512},
		Timeout:         12 * 60 * 60, // 12 hours
	}
	result.Use(options...)
	return result
}

// Verify verifies the given http.Request. This is
// syntactic sugar for NewVerifier(options...).Verify(request)
func Verify(request *http.Request, body []byte, certificate string, options ...VerifierOption) error {
	verifier := NewVerifier(options...)
	return verifier.Verify(request, body, certificate)
}

// Use applies the given options to the Verifier
func (verifier *Verifier) Use(options ...VerifierOption) {
	for _, option := range options {
		option(verifier)
	}
}

// Verify verifies the given http.Request
func (verifier *Verifier) Verify(request *http.Request, body []byte, certificate string) error {

	const location = "hannibal.sigs.Verify"

	if request == nil {
		return derp.NewInternalError("hannibal.sigs.Verify", "Request cannot be nil")
	}

	log.Debug().
		Str("certificate", certificate).
		Msg("Verifying Signature")

	// Verify the body Digest
	if err := VerifyDigest(request, body, verifier.BodyDigests...); err != nil {
		return derp.Wrap(err, location, "Error verifying body digest")
	}

	// Retrieve and parse the Signature from the HTTP Request
	signature, err := ParseSignature(request.Header.Get("Signature"))

	if err != nil {
		return derp.Wrap(err, location, "Error parsing signature")
	}

	log.Trace().
		Interface("signature", signature).
		Msg("Parsed Signature")

	// RULE: Verify that the signature contains all of the fields that we require
	if !slice.ContainsAll(signature.Headers, verifier.Fields...) {
		return derp.NewForbiddenError("hannibal.sigs.Verify", "Signature must include ALL of these fields", verifier.Fields)
	}

	// Decode the PEM certificate into a public key
	publicKey, err := DecodePublicPEM(certificate)

	if err != nil {
		return derp.Wrap(err, location, "Error decoding public key")
	}

	// Recreate the plaintext and digest used to make the Signature
	plaintext := makePlaintext(request, signature.Headers...)

	// Try each hash in order
	for _, hash := range verifier.SignatureHashes {
		if err := verifyHashAndSignature(plaintext, hash, publicKey, signature.Signature); err == nil {
			log.Trace().Msg("Trying " + hash + ": Succeeded")
			return nil
		} else {
			log.Trace().Err(err).Msg("Trying " + hash + ": Failed")
		}
	}

	return derp.NewForbiddenError(location, "Invalid signature")
}

/******************************************
 * Helper Functions
 ******************************************/

// Verify Hash And Signature computes the hashed value of the plaintext, then verifies
// that this result matches the provided public key and signature.  It returns an error
// if the signature does not match.
func verifyHashAndSignature(plaintext string, hash string, publicKey crypto.PublicKey, signature []byte) error {

	const location = "hannibal.sigs.verifyHashAndSignature"

	// Make a digest using the hash algorithm
	digest, err := makeSignatureHash(plaintext, hash)

	if err != nil {
		return derp.Wrap(err, location, "Error creating digest")
	}

	if log.Logger.GetLevel() == zerolog.TraceLevel {
		log.Trace().
			Str("plaintext", plaintext).
			Str("hash", hash).
			Str("signature", base64.StdEncoding.EncodeToString(signature)).
			Msg("VerifyHashAndSignature")
	}

	// Verify the signature matches the message digest
	if err := verifySignature(digest, signature, publicKey); err != nil {
		return derp.Wrap(err, location, "Invalid signature")
	}

	// Beauty is in the eye of the beholder.
	return nil
}

// verifySignature verifies the given signature using the provided public key.
// The public key can be either an RSA or ECDSA keys.
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
