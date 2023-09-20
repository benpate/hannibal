package sigs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/slice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Verifier contains all of the settings necessary to verify a request
type Verifier struct {
	Fields          []string
	BodyDigests     []crypto.Hash // List of algorithms to accept from remote servers when they create a Digest header.  Default is SHA256 and SHA512
	SignatureHashes []crypto.Hash // Digest algorithm used to create the signature.  Default is SHA256, SHA512
	Timeout         int           // Number of seconds before signatures are expired. Default is 43200 seconds (12 hours).
	CheckDigest     bool          // If true, then the verifier will check the Digest header.  Default is true.
}

// NewVerifier returns a fully initialized Verifier
func NewVerifier(options ...VerifierOption) Verifier {
	result := Verifier{
		Fields:          []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		BodyDigests:     []crypto.Hash{crypto.SHA256, crypto.SHA512},
		SignatureHashes: []crypto.Hash{crypto.SHA256, crypto.SHA512},
		Timeout:         12 * 60 * 60, // 12 hours
		CheckDigest:     true,
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

	if request == nil {
		return derp.NewInternalError("hannibal.sigs.Verify", "Request cannot be nil")
	}

	log.Debug().
		Str("loc", location).
		Msg("Verifying Signature")

	// Verify the request date
	if verifier.Timeout > 0 {
		date, err := time.Parse(http.TimeFormat, request.Header.Get(FieldDate))

		if err != nil {
			return derp.Wrap(err, location, "Invalid Date header.  Must match 'Mon, 02 Jan 2006 15:04:05 GMT'")
		}

		if date.Unix() < time.Now().Add(-1*time.Duration(verifier.Timeout)*time.Second).Unix() {
			return derp.NewForbiddenError(location, "Request date has expired. Must be within the last "+strconv.Itoa(verifier.Timeout)+" seconds")
		}
	}

	// Verify the body Digest (default behavior)
	if verifier.CheckDigest {
		if err := VerifyDigest(request, verifier.BodyDigests...); err != nil {
			return derp.Wrap(err, location, "Error verifying body digest")
		}
	}

	// Retrieve and parse the Signature from the HTTP Request
	signature, err := ParseSignature(request.Header.Get("Signature"))

	if err != nil {
		return derp.Wrap(err, location, "Error parsing signature")
	}

	log.Trace().
		Str("loc", location).
		Str("certificate", certificate).
		Interface("signature", signature).
		Msg("Parsed Signature")

	// RULE: If the signature has expired, then reject it.
	if signature.IsExpired(verifier.Timeout) {
		return derp.NewForbiddenError(location, "Signature has expired")
	}

	// RULE: Verify that the signature contains all of the fields that we require
	if !slice.ContainsAll(signature.Headers, verifier.Fields...) {
		return derp.NewForbiddenError(location, "Signature must include ALL of these fields", verifier.Fields)
	}

	// Decode the PEM certificate into a public key
	publicKey, err := DecodePublicPEM(certificate)

	if err != nil {
		return derp.Wrap(err, location, "Error decoding public key")
	}

	// Recreate the plaintext and digest used to make the Signature
	plaintext := makePlaintext(request, signature, signature.Headers...)

	// Try each hash in order
	for _, hash := range verifier.SignatureHashes {
		if err := verifyHashAndSignature(plaintext, hash, publicKey, signature.Signature); err == nil {
			log.Debug().Str("loc", location).Msg("Signature is VALID")
			return nil
		}
	}

	return derp.NewForbiddenError(location, "Signature is INVALID")
}

/******************************************
 * Helper Functions
 ******************************************/

// Verify Hash And Signature computes the hashed value of the plaintext, then verifies
// that this result matches the provided public key and signature.  It returns an error
// if the signature does not match.
func verifyHashAndSignature(plaintext string, hash crypto.Hash, publicKey crypto.PublicKey, signature []byte) error {

	const location = "hannibal.sigs.verifyHashAndSignature"

	// Make a digest using the hash algorithm
	digest, err := makeSignatureHash(plaintext, hash)

	if err != nil {
		return derp.Wrap(err, location, "Error creating digest")
	}

	// Logging here.. wrapping it in an "if" because the base64 encoding is expensive
	if log.Logger.GetLevel() == zerolog.TraceLevel {
		log.Trace().
			Str("plaintext", plaintext).
			Int("hash", int(hash)).
			Str("signature", base64.StdEncoding.EncodeToString(signature)).
			Str("digest", base64.StdEncoding.EncodeToString(digest)).
			Msg("VerifyHashAndSignature")
	}

	// Verify the signature matches the message digest
	if err := verifySignature(publicKey, hash, digest, signature); err != nil {
		err = derp.Wrap(err, location, "Invalid signature")
		log.Debug().Err(err).Msg("Signature is Invalid")
		return err
	}

	// Beauty is in the eye of the beholder.
	return nil
}

// verifySignature verifies the given signature using the provided public key.
// The public key can be either an RSA or ECDSA keys.
func verifySignature(publicKey crypto.PublicKey, hash crypto.Hash, digest []byte, signature []byte) error {

	const location = "hannibal.sigs.verifySignature"

	switch typedKey := publicKey.(type) {

	case *rsa.PublicKey:
		if err := rsa.VerifyPKCS1v15(typedKey, hash, digest, signature); err != nil {
			return derp.Wrap(err, location, "Error verifying RSA signature")
		}
		return nil

	case *ecdsa.PublicKey:
		if !ecdsa.VerifyASN1(typedKey, digest, signature) {
			return derp.NewForbiddenError(location, "Invalid ECDSA signature")
		}
		return nil
	}

	return derp.NewBadRequestError("hannibal.sigs.verifySignature", "Unrecognized public key type")
}
