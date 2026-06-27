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

// getDigestFuncByName returns the DigestFunc for either `sha-256` or `sha-512`.
// Unrecognized names return an error so callers can skip them, rather than
// silently falling back to SHA-256 (which getHashByName does).
func getDigestFuncByName(name string) (DigestFunc, error) {

	hash, ok := lookupHashByName(name)

	if !ok {
		return nil, derp.BadRequest("sigs.getDigestFuncByName", "Unknown digest algorithm", name)
	}

	return getDigestFunc(hash)
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

	return nil, derp.BadRequest("sigs.getDigestFunc", "Unknown algorithm", algorithm)
}

// getDigestName returns the name of a given crypto.Hash value
func getDigestName(algorithm crypto.Hash) string {

	switch algorithm {

	case crypto.SHA256:
		return "SHA-256"

	case crypto.SHA512:
		return "SHA-512"
	}

	return "unknown"
}

// getHashByName converts common hash names into crypto.Hash values.  It works
// with these values: sha-256, sha256, sha-512, sha512 (case insensitive).
// Unknown names default to SHA-256; callers that must distinguish an unknown
// name should use lookupHashByName instead.
func getHashByName(name string) crypto.Hash {

	if hash, ok := lookupHashByName(name); ok {
		return hash
	}

	log.Warn().Msg("sigs.getHashByName: Unknown hash name: " + name + ". Defaulting to SHA-256")

	return crypto.SHA256
}

// lookupHashByName converts common hash names into crypto.Hash values, reporting
// whether the name was recognized. It accepts sha-256, sha256, sha-512, sha512
// (case insensitive).
func lookupHashByName(name string) (crypto.Hash, bool) {

	switch strings.ToLower(name) {

	case "sha-256", "sha256":
		return crypto.SHA256, true

	case "sha-512", "sha512":
		return crypto.SHA512, true
	}

	return 0, false
}

// DigestSHA256 calculates the SHA-256 digest of a slice of bytes
func DigestSHA256(body []byte) string {
	digest := sha256.Sum256(body)
	return base64.StdEncoding.EncodeToString(digest[:])
}

// DigestSHA512 calculates the SHA-512 digest of a given slice of bytes
func DigestSHA512(body []byte) string {
	digest := sha512.Sum512(body)
	return base64.StdEncoding.EncodeToString(digest[:])
}

// Additional digest algorithms are not yet supported: unixsum, unixcksum, crc32c,
// sha-256 and sha-512, id-sha-256, id-sha-512.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Digest
