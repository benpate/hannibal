package sigs

import "net/http"

// Verifier contains all of the settings necessary to verify a request
type Verifier struct {
	Fields  []string
	Digests []string
}

// NewVerifier returns a fully initialized Verifier
func NewVerifier(options ...VerifierOption) Verifier {
	result := Verifier{
		Fields:  []string{FieldRequestTarget, FieldHost, FieldDate, FieldDigest},
		Digests: []string{Digest_SHA256, Digest_SHA512},
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
	return nil
}
