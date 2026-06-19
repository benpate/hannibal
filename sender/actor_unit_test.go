package sender

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewActor confirms the default Actor implementation reports the ID and
// private-key material it was constructed with.
func TestNewActor(t *testing.T) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	actor := NewActor(
		"https://example.com/users/alice",
		"https://example.com/users/alice#main-key",
		privateKey,
	)

	assert.Equal(t, "https://example.com/users/alice", actor.ActorID())

	keyID, key := actor.PrivateKey()
	assert.Equal(t, "https://example.com/users/alice#main-key", keyID)
	assert.Equal(t, privateKey, key)
}
