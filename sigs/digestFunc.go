package sigs

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"strings"

	"github.com/benpate/derp"
)

// DigestFunc defines a function that calculates the digest of a given byte array
type DigestFunc func(body []byte) string

// getDigestFuncs uses a list of algorithm names to generate a list of DigestFuncs
// nolint:unused // We may use this later, so just keep it for nao.
func getDigestFuncs(algorithms ...string) ([]DigestFunc, error) {

	result := make([]DigestFunc, 0, len(algorithms))

	for _, algorithm := range algorithms {

		fn, err := getDigestFunc(algorithm)

		if err != nil {
			return nil, derp.Wrap(err, "sigs.getDigestFunc", "Error parsing algorithm")
		}

		result = append(result, fn)
	}

	return result, nil
}

// getDigestFunc uses an algorithm name to generate a DigestFunc using
// a case insensitive match.  It currently supports `sha-256` and `sha-512`.
// Unrecognized digest names will return an error.
func getDigestFunc(algorithm string) (DigestFunc, error) {

	switch strings.ToLower(algorithm) {

	case Digest_SHA256:
		return DigestSHA256, nil

	case Digest_SHA512:
		return DigestSHA512, nil
	}

	return nil, derp.NewBadRequestError("sigs.getDigestFunc", "Unknown algorithm: %s", algorithm)
}

// DigestSHA256 calculates the SHA-256 digest of a slice of bytes
func DigestSHA256(body []byte) string {
	digest := sha256.Sum256(body)
	return "SHA-256=" + base64.StdEncoding.EncodeToString(digest[:])
}

// DigestSHA512 calculates the SHA-512 digest of a given slice of bytes
func DigestSHA512(body []byte) string {
	digest := sha512.Sum512(body)
	return "SHA-512=" + base64.StdEncoding.EncodeToString(digest[:])
}
