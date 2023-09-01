package sigs

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcDigest(t *testing.T) {

	var body bytes.Buffer

	body.WriteString("This is my body. There are many like it, but this one is mine.")

	request, err := http.NewRequest("GET", "http://example.com/foo", &body)
	require.Nil(t, err)

	digest, err := CalcDigest(request, DigestSHA256)
	require.Nil(t, err)
	require.Equal(t, "SHA-256=2dZxOmbiuR4yypVcyCfajB3YMhmSg+QNUlnUIrfllPM=", digest)
}

func TestApplyDigest(t *testing.T) {

	var body bytes.Buffer

	body.WriteString("This is my body. There are many like it, but this one is mine.")

	request, err := http.NewRequest("GET", "http://example.com/foo", &body)
	require.Nil(t, err)

	err = ApplyDigest(request, DigestSHA256)
	require.Nil(t, err)
	require.Equal(t, "SHA-256=2dZxOmbiuR4yypVcyCfajB3YMhmSg+QNUlnUIrfllPM=", request.Header.Get("Digest"))
}

func TestVerifyDigest(t *testing.T) {

	var body bytes.Buffer

	body.WriteString("This is my body. There are many like it, but this one is mine.")

	request, err := http.NewRequest("GET", "http://example.com/foo", &body)
	require.Nil(t, err)

	err = ApplyDigest(request, DigestSHA256)
	require.Nil(t, err)

	err = VerifyDigest(request)
	require.Nil(t, err)
}
