package sigs

import (
	"crypto"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMakeSignedDigest_UnrecognizedKey confirms signing with a key that is
// neither RSA nor ECDSA returns an error rather than producing a bogus signature.
func TestMakeSignedDigest_UnrecognizedKey(t *testing.T) {

	digest := []byte("0123456789abcdef0123456789abcdef")

	result, err := makeSignedDigest(digest, crypto.SHA256, "not-a-real-key")
	require.NotNil(t, err, "an unrecognized private key type must be rejected")
	require.Nil(t, result)
}
