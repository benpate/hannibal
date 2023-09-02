package sigs

// SignerOption is a function that modifies a Signer
type SignerOption func(*Signer)

// SignerFields sets the http.Request fields to be signed
func SignerFields(fields ...string) SignerOption {
	return func(signer *Signer) {
		signer.Fields = fields
	}
}

// SignerSignatureDigest sets the hashing algorithm to be used
// when we sign a request.
func SignerSignatureHash(hash string) SignerOption {
	return func(signer *Signer) {
		signer.SignatureHash = hash
	}
}

// SignerBodyDigests sets the digest algorithm to be used
// when creating the "Digest" header.
func SignerBodyDigest(digest string) SignerOption {
	return func(signer *Signer) {
		signer.BodyDigest = digest
	}
}
