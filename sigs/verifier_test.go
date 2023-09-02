package sigs

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {

	// Configure logging
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// Create a new Request
	bodyReader := strings.NewReader("This is the body of the request")

	request, err := http.NewRequest("GET", "http://example.com/something?test=true", bodyReader)
	require.Nil(t, err)
	request.Header.Set("Content-Type", "text/plain")

	// Create a Private Key to sign the request
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// Sign the Request
	err = Sign(request, "test-key", privateKey)
	require.Nil(t, err)

	require.Equal(t, "SHA-256=65F8+S1Bg7oPQS/fIxVg4x7PoLWnOxWlGMFB/hafojg=", request.Header.Get("Digest"))
	require.NotEmpty(t, request.Header.Get("Signature"))

	// Verify the Request
	publicKeyPEM := EncodePublicPEM(privateKey)

	err = Verify(request, publicKeyPEM)
	require.Nil(t, err)
}

func TestSignAndVerify(t *testing.T) {

	// Create an RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// Create a test digest
	digest := []byte("test")

	// Sign the digest
	signature, err := makeSignedDigest(digest, privateKey)
	require.Nil(t, err)

	err = verifySignature(digest, signature, &privateKey.PublicKey)
	require.Nil(t, err)
}
