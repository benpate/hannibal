package sigs

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignRequest(t *testing.T) {

	bodyReader := strings.NewReader("This is the body of the request")

	request, err := http.NewRequest("GET", "http://example.com/something?test=true", bodyReader)
	require.Nil(t, err)
	request.Header.Set("Content-Type", "text/plain")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	err = Sign(request, "test-key", privateKey)
	require.Nil(t, err)

	require.Equal(t, "SHA-256=65F8+S1Bg7oPQS/fIxVg4x7PoLWnOxWlGMFB/hafojg=", request.Header.Get("Digest"))
	require.NotEmpty(t, request.Header.Get("Signature"))
}

func TestMakePlaintext(t *testing.T) {

	bodyReader := strings.NewReader("This is the body of the request")

	request, err := http.NewRequest("GET", "http://example.com/something?test=true", bodyReader)
	require.Nil(t, err)
	request.Header.Set("Content-Type", "text/plain")

	result := makePlaintext(request, FieldRequestTarget, FieldHost, "Content-Type")
	expected := removeTabs(
		`(request-target): get /something?test=true
		host: example.com
		content-type: text/plain`)

	require.Equal(t, expected, result)
}

func TestMakeDigest_SHA256(t *testing.T) {
	result, err := makePlaintextDigest("This is digest-able", "sha-256")
	require.Nil(t, err)

	actual := base64.StdEncoding.EncodeToString(result)
	require.Equal(t, "jlBmJDmZdMjhLZga/ZjDrlloKd5lukG9S0lu/f7Xc64=", actual)
}

func TestMakeDigest_SHA512(t *testing.T) {
	result, err := makePlaintextDigest("This is digest-able", "sha-512")
	require.Nil(t, err)

	actual := base64.StdEncoding.EncodeToString(result)
	require.Equal(t, "s2JJ/rYbVQTrkNR440jq+wuNk9ktJgvmVSDq805iC0EP4ONQPwfvuQK0yR/YuX7riJtNRwxMq6R1GL8W7A5vzg==", actual)
}

func TestMakeDigest_Other(t *testing.T) {
	result, err := makePlaintextDigest("This is digest-able", "invalid-hash-function")
	require.NotNil(t, err)
	require.Nil(t, result)
}

func TestGetPathAndQuery(t *testing.T) {
	url, _ := url.Parse("http://example.com")
	require.Equal(t, "/", getPathAndQuery(url))

	url, _ = url.Parse("http://example.com/")
	require.Equal(t, "/", getPathAndQuery(url))

	url, _ = url.Parse("http://example.com/something")
	require.Equal(t, "/something", getPathAndQuery(url))

	url, _ = url.Parse("http://example.com/something?test=true")
	require.Equal(t, "/something?test=true", getPathAndQuery(url))
}

func removeTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "")
}
