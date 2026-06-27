package streams

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// allowPrivateIPs is a remote.Option that lets a transaction connect to the
// loopback address, which the SSRF guard blocks by default. It is needed because
// the test server below listens on 127.0.0.1.
func allowPrivateIPs() remote.Option {
	return remote.Option{
		BeforeRequest: func(txn *remote.Transaction) error {
			txn.AllowPrivateIPs(true)
			return nil
		},
	}
}

// TestDefaultClient_Load confirms the default client fetches and parses a remote
// ActivityPub document.
func TestDefaultClient_Load(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", vocab.ContentTypeActivityPub)
		_, _ = w.Write([]byte(`{"id":"urn:loaded","type":"Note"}`))
	}))
	defer server.Close()

	client := NewDefaultClient(allowPrivateIPs())
	document, err := client.Load(server.URL)

	require.NoError(t, err)
	assert.Equal(t, "urn:loaded", document.ID())
	assert.Equal(t, vocab.ObjectTypeNote, document.Type())

	// The response header is attached to the returned document.
	assert.Equal(t, vocab.ContentTypeActivityPub, document.HTTPHeader().Get("Content-Type"))
}

// TestDefaultClient_Load_Error confirms a server error is surfaced to the caller.
func TestDefaultClient_Load_Error(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewDefaultClient(allowPrivateIPs())
	_, err := client.Load(server.URL)

	require.Error(t, err)
}

// TestDefaultClient_SaveDeleteNoop confirms Save and Delete are no-ops that
// return nil, and SetRootClient does not panic.
func TestDefaultClient_SaveDeleteNoop(t *testing.T) {

	client := NewDefaultClient()
	document := NewDocument(map[string]any{vocab.PropertyID: "urn:test"})

	assert.NoError(t, client.Save(document))
	assert.NoError(t, client.Delete(document.ID()))

	// SetRootClient is a no-op for the default client; calling it must not panic.
	assert.NotPanics(t, func() { client.SetRootClient(client) })
}
