package sigs

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyRSA_SHA(t *testing.T) {

	var plaintext = []byte("This is the message to be signed and verified")

	sha256 := crypto.SHA256.New()
	sha256.Write(plaintext)
	hashed := sha256.Sum(nil)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	{
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
		require.Nil(t, err)
		fmt.Println(base64.StdEncoding.EncodeToString(signature))
	}
	fmt.Println("---------")
	{
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, plaintext)
		require.Nil(t, err)
		fmt.Println(base64.StdEncoding.EncodeToString(signature))
	}
}
