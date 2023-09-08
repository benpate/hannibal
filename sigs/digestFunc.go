package sigs

import (
	"crypto"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"strings"

	"github.com/benpate/derp"
	"github.com/rs/zerolog/log"
)

// DigestFunc defines a function that calculates the digest of a given byte array
type DigestFunc func(body []byte) string

// getDigestFuncs uses a list of algorithm names to generate a list of DigestFuncs
// nolint:unused // We may use this later, so just keep it for nao.
func getDigestFuncs(algorithms ...crypto.Hash) ([]DigestFunc, error) {

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

// getDigestFuncByName returns the DigestFunc for either `sha-256` or `sha-512`
func getDigestFuncByName(name string) (DigestFunc, error) {
	return getDigestFunc(getHashByName(name))
}

// getDigestFunc uses an algorithm name to generate a DigestFunc using
// a case insensitive match.  It currently supports `sha-256` and `sha-512`.
// Unrecognized digest names will return an error.
func getDigestFunc(algorithm crypto.Hash) (DigestFunc, error) {

	switch algorithm {

	case crypto.SHA256:
		return DigestSHA256, nil

	case crypto.SHA512:
		return DigestSHA512, nil
	}

	return nil, derp.NewBadRequestError("sigs.getDigestFunc", "Unknown algorithm", algorithm)
}

// getHashByName converts common hash names into crypto.Hash values.  It works
// with these values: sha-256, sha256, sha-512, sha512 (case insensitive)
func getHashByName(name string) crypto.Hash {

	switch strings.ToLower(name) {

	case "sha-256", "sha256":
		return crypto.SHA256

	case "sha-512", "sha512":
		return crypto.SHA512
	}

	log.Warn().Msg("sigs.getHashByName: Unknown hash name: " + name + ". Defaulting to SHA-256")

	return crypto.SHA256
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

// TODO: Additional algorithms specified by https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Digest
// unixsum, unixcksum, crc32c, sha-256 and sha-512, id-sha-256, id-sha-512
