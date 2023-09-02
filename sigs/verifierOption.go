package sigs

// VerifierOption is a function that modifies a Verifier
type VerifierOption func(*Verifier)

// VerifierFields sets the list of http.Request fields that
// MUST ALL be present in the "Signature" header from a
// remote server for a signature to be accepted. Extra
// fields are allowed in the Signature, and will still
// be verified.
func VerifierFields(fields ...string) VerifierOption {
	return func(verifier *Verifier) {
		verifier.Fields = fields
	}
}

// VerifierDigests sets the list of algorithms that we will
// accept from remote servers when they create a "Digest"
// http header. ALL recognized digests must be valid to
// pass, and AT LEAST ONE of the algorithms must be from
// this list.
func VerifierBodyDigests(digests ...string) VerifierOption {
	return func(verifier *Verifier) {
		verifier.BodyDigests = digests
	}
}

// VerifierSignatureHashes sets the hashing algorithms to use
// when validating the "Signature" header. Hashes are tried
// in order, and the FIRST successful match returns success.
// If ALL hash attempts fail, then validation fails.
func VerifierSignatureHashes(hashes ...string) VerifierOption {
	return func(verifier *Verifier) {
		verifier.SignatureHashes = hashes
	}
}
