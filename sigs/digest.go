package sigs

import (
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
	"github.com/rs/zerolog/log"
)

// ApplyDigest calculates the digest of the body from a given
// http.Request, then adds the digest to the Request's header.
func ApplyDigest(request *http.Request, body []byte, fn DigestFunc) error {

	// Try to calculate the digest with the DigestFunc
	result := fn(body)

	// Apply the digest to the Request
	request.Header.Set(FieldDigest, result)
	return nil
}

// VerifyDigest verifies that the digest in the http.Request header
// matches the contents of the http.Request body.
func VerifyDigest(request *http.Request, body []byte, allowedAlgorithms ...string) error {

	// Retrieve the digest(s) included in the HTTP Request
	digestHeader := request.Header.Get(FieldDigest)

	// If there is no digest header, then there is nothing to verify
	if digestHeader == "" {
		return nil
	}

	// Process the digest header into separate values
	headerValues := strings.Split(digestHeader, ",")
	atLeastOneAlgorithmMatches := false

	for _, headerValue := range headerValues {

		headerValue = strings.Trim(headerValue, " ")
		digestAlgorithm, digestValue := list.Split(headerValue, '=')

		// If we recognize the digest algorithm, then use it to verify the body/digest
		fn, err := getDigestFunc(digestAlgorithm)

		if err != nil {
			log.Trace().Msg("sigs.VerifyDigest: Unknown digest algorithm: " + digestAlgorithm)
			continue
		}

		// If the values match, then success!
		if headerValue == fn(body) {
			log.Trace().Msg("sigs.VerifyDigest: Valid Digest Found. Algorithm: " + digestAlgorithm)
			atLeastOneAlgorithmMatches = true
			continue
		}

		// If the values DON'T MATCH, then fail immediately.
		// We don't want bad actors "digest shopping"
		return derp.NewForbiddenError("sigs.VerifyDigest", "Digest verification failed", digestValue)
	}

	// If we have found at least one digest that matches, then success!
	if atLeastOneAlgorithmMatches {
		return nil
	}

	// Otherwise, the digest hash does not meet our minimum requirements.  Fail.
	return derp.NewForbiddenError("sigs.VerifyDigest", "No matching digest found")
}
