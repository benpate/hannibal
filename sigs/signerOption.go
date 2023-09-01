package sigs

// SignerOption is a function that modifies a Signer
type SignerOption func(*Signer)

// SignFields sets the http.Request fields to be signed
func SignFields(fields ...string) SignerOption {
	return func(signer *Signer) {
		signer.Fields = fields
	}
}

// SignSignatureDigest sets the digest algorithm to be used when
// signing the request.
func SignSignatureDigest(digest string) SignerOption {
	return func(signer *Signer) {
		signer.SignatureDigest = digest
	}
}

// SignBodyDigests sets the digest algorithms to be used when creating
// the "Digest" header.
func SignBodyDigests(digests ...string) SignerOption {
	return func(signer *Signer) {
		signer.BodyDigests = digests
	}
}
