package sigs

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/list"
	"github.com/rs/zerolog/log"
)

// CalcDigest uses a DigestFunc to calculate the digest from the body
// of a given http.Request.
func CalcDigest(request *http.Request, fn DigestFunc) (string, error) {

	var body bytes.Buffer

	// Try to get a copy of the Request body
	bodyReader, err := request.GetBody()

	if err != nil {
		return "", derp.Wrap(err, "pub.RequestDigest", "Error getting request body")
	}

	// Try to read the request body into a buffer
	if _, err := body.ReadFrom(bodyReader); err != nil {
		return "", derp.Wrap(err, "pub.RequestDigest", "Error reading request body")
	}

	// Calculate the digest with the DigestFunc
	return fn(body.Bytes()), nil
}

// ApplyDigest calculates the digest of the body from a given
// http.Request, then adds the digest to the Request's header.
func ApplyDigest(request *http.Request, fn DigestFunc) error {

	// Try to calculate the digest with the DigestFunc
	result, err := CalcDigest(request, fn)

	if err != nil {
		return derp.Wrap(err, "pub.RequestDigest", "Error calculating digest")
	}

	// Apply the digest to the Request
	request.Header.Set(FieldDigest, result)
	return nil
}

// VerifyDigest verifies that the digest in the http.Request header
// matches the contents of the http.Request body.
func VerifyDigest(request *http.Request) error {

	var body bytes.Buffer

	// Try to get a copy of the Request body
	bodyReader, err := request.GetBody()

	if err != nil {
		return derp.Wrap(err, "pub.RequestDigest", "Error getting request body")
	}

	// Try to read the request body into a buffer
	if _, err := body.ReadFrom(bodyReader); err != nil {
		return derp.Wrap(err, "pub.RequestDigest", "Error reading request body")
	}

	// Retrieve the digest(s) included in the HTTP Request
	header := request.Header.Get(FieldDigest)
	headerValues := strings.Split(header, ",")

	// Scan multiple digest values
	for _, headerValue := range headerValues {

		digestAlgorithm, digestValue := list.Split(headerValue, '=')

		// If we recognize the digest algorithm, then use it to verify the body/digest
		fn, err := getDigestFunc(digestAlgorithm)

		if err != nil {
			log.Trace().Msg("sigs.VerifyDigest: Unknown digest algorithm: " + digestAlgorithm)
			continue
		}

		// If the values match, then success!
		if headerValue == fn(body.Bytes()) {
			log.Trace().Msg("sigs.VerifyDigest: Valid Digest Found. Algorithm: " + digestAlgorithm)
			return nil
		}

		// If the values DON'T MATCH, then fail immediately.
		// We don't want bad actors "digest shopping"
		return derp.NewForbiddenError("sigs.VerifyDigest", "Digest verification failed", digestValue)
	}

	// Fall through means that we could not find a digest that we can use
	return derp.NewForbiddenError("sigs.VerifyDigest", "No matching digest found")
}
