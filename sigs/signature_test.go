package sigs

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

var testString string
var testBytes []byte

func init() {
	testString = "Y2FiYWIxNGRiZDk4ZA=="
	testBytes, _ = base64.StdEncoding.DecodeString(testString)
}

func TestParseSignature_IETF(t *testing.T) {

	// From: https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures-05#section-4.1.1
	headerValue := `keyId="rsa-key-1",algorithm="hs2019",
	headers="(request-target) (created) host digest content-length",
	signature="` + testString + `"`

	signature, err := ParseSignature(headerValue)

	require.Nil(t, err)
	require.Equal(t, "rsa-key-1", signature.KeyID)
	require.Equal(t, "hs2019", signature.Algorithm)
	require.Equal(t, []string{"(request-target)", "(created)", "host", "digest", "content-length"}, signature.Headers)
	require.Equal(t, testBytes, signature.Signature)
}

func TestParseSignature_Mastodon(t *testing.T) {

	// From: https://docs.joinmastodon.org/spec/security/
	headerValue := `keyId="https://my-example.com/actor#main-key",headers="(request-target) host date digest",signature="` + testString + `"`

	signature, err := ParseSignature(headerValue)

	require.Nil(t, err)
	require.Equal(t, "https://my-example.com/actor#main-key", signature.KeyID)
	require.Equal(t, "", signature.Algorithm)
	require.Equal(t, []string{"(request-target)", "host", "date", "digest"}, signature.Headers)
	require.Equal(t, testBytes, signature.Signature)
}

func TestParseSignature_MastodonWithWhitespace(t *testing.T) {

	// From: https://docs.joinmastodon.org/spec/security/
	// Adding spaces to the signature just in case others do this too.
	headerValue := `keyId="https://my-example.com/actor#main-key", headers="(request-target) host date digest", signature="` + testString + `"`

	signature, err := ParseSignature(headerValue)

	require.Nil(t, err)
	require.Equal(t, "https://my-example.com/actor#main-key", signature.KeyID)
	require.Equal(t, "", signature.Algorithm)
	require.Equal(t, []string{"(request-target)", "host", "date", "digest"}, signature.Headers)
	require.Equal(t, testBytes, signature.Signature)
}
