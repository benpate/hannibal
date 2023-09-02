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

// VerifierSignatureHash sets the hashing algorithm to use
// when validating the "Signature" header. According to
// the http signatures spec, this should always be sha-256.
func VerifierSignatureHash(hash string) VerifierOption {
	return func(verifier *Verifier) {
		verifier.SignatureHash = hash
	}
}
