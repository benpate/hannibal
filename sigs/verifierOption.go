package sigs

// VerifierOption is a function that modifies a Verifier
type VerifierOption func(*Verifier)

// VerifyFields sets the http.Request fields to be signed
func VerifyFields(fields ...string) VerifierOption {
	return func(verifier *Verifier) {
		verifier.Fields = fields
	}
}

// VerifyDigests sets the algorithms to be used when creating
// the "Digest" header.
func VerifyDigests(digests ...string) VerifierOption {
	return func(verifier *Verifier) {
		verifier.Digests = digests
	}
}
