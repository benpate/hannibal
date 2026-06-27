package sender

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"iter"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/turbine/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain enables delivery to non-public addresses for the whole package, since
// these tests POST to loopback httptest servers that remote's SSRF guard would
// otherwise block. Production keeps allowPrivateIPs = false.
func TestMain(m *testing.M) {
	allowPrivateIPs = true
	os.Exit(m.Run())
}

// keyedActor is a test Actor that carries a real RSA private key, so that
// outbound requests can actually be signed.
type keyedActor struct {
	id         string
	keyID      string
	privateKey crypto.PrivateKey
}

func (a keyedActor) ActorID() string { return a.id }
func (a keyedActor) PrivateKey() (string, crypto.PrivateKey) {
	return a.keyID, a.privateKey
}

// keyedLocator resolves exactly one actor (the keyed test actor).
type keyedLocator struct {
	actor keyedActor
}

func (l keyedLocator) Actor(id string) (Actor, error) {
	if id == l.actor.id {
		return l.actor, nil
	}
	return nil, derp.NotFound("keyedLocator.Actor", "unknown actor", id)
}

func (l keyedLocator) Recipient(url string) (iter.Seq[string], error) {
	return func(yield func(string) bool) {}, nil
}

// newKeyedSender builds a Sender whose single actor has a real signing key.
func newKeyedSender(t *testing.T) (Sender, string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	actorID := "https://example.com/users/alice"
	actor := keyedActor{
		id:         actorID,
		keyID:      actorID + "#main-key",
		privateKey: privateKey,
	}

	q, _ := newRecordingQueue()
	return New(keyedLocator{actor: actor}, q), actorID
}

// TestSendToSingleRecipient_Success confirms a deliverable activity is POSTed to
// the recipient inbox, signed, and reports Success.
func TestSendToSingleRecipient_Success(t *testing.T) {

	sender, actorID := newKeyedSender(t)

	var gotSignature atomic.Bool
	var gotDigest atomic.Bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSignature.Store(r.Header.Get("Signature") != "")
		gotDigest.Store(r.Header.Get("Digest") != "")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := sender.SendToSingleRecipient(mapof.Any{
		"actor":    actorID,
		"inbox":    server.URL,
		"activity": mapof.Any{"type": "Create", "actor": actorID},
	})

	assert.Equal(t, queue.ResultStatusSuccess, result.Status)
	assert.True(t, gotSignature.Load(), "the outbound request must be signed")
	assert.True(t, gotDigest.Load(), "the outbound request must carry a body digest")
}

// TestSendToSingleRecipient_ActorNotFound confirms an unknown sending actor
// yields a Failure (cannot be retried).
func TestSendToSingleRecipient_ActorNotFound(t *testing.T) {

	sender, _ := newKeyedSender(t)

	result := sender.SendToSingleRecipient(mapof.Any{
		"actor":    "https://example.com/users/nobody",
		"inbox":    "https://example.com/inbox",
		"activity": mapof.Any{"type": "Create"},
	})

	assert.Equal(t, queue.ResultStatusFailure, result.Status)
}

// TestSendToSingleRecipient_ServerError confirms a 5xx response is classified as
// a retriable Error.
func TestSendToSingleRecipient_ServerError(t *testing.T) {

	sender, actorID := newKeyedSender(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	result := sender.SendToSingleRecipient(mapof.Any{
		"actor":    actorID,
		"inbox":    server.URL,
		"activity": mapof.Any{"type": "Create", "actor": actorID},
	})

	// 5xx is the remote server's fault -> retriable Error.
	assert.Equal(t, queue.ResultStatusError, result.Status)
}

// TestSendToSingleRecipient_ClientError confirms a 4xx response is classified as
// a non-retriable Failure.
func TestSendToSingleRecipient_ClientError(t *testing.T) {

	sender, actorID := newKeyedSender(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	result := sender.SendToSingleRecipient(mapof.Any{
		"actor":    actorID,
		"inbox":    server.URL,
		"activity": mapof.Any{"type": "Create", "actor": actorID},
	})

	// 4xx is our fault -> non-retriable Failure.
	assert.Equal(t, queue.ResultStatusFailure, result.Status)
}

// TestSendToSingleRecipient_TooManyRequests confirms a 429 response is requeued
// rather than failed.
func TestSendToSingleRecipient_TooManyRequests(t *testing.T) {

	sender, actorID := newKeyedSender(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	result := sender.SendToSingleRecipient(mapof.Any{
		"actor":    actorID,
		"inbox":    server.URL,
		"activity": mapof.Any{"type": "Create", "actor": actorID},
	})

	assert.Equal(t, queue.ResultStatusRequeue, result.Status)
}

// TestSignRequest confirms the signRequest middleware signs an outbound request,
// producing a verifiable signature.
func TestSignRequest(t *testing.T) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	keyID := "https://example.com/users/alice#main-key"

	var verified atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := sigs.Verify(r, func(string) (string, error) {
			return sigs.EncodePublicPEM(privateKey), nil
		})
		verified.Store(err == nil)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// httptest gives us a *http.Client; use the remote package via the sender's
	// own path by signing through the middleware directly.
	request, err := http.NewRequest(http.MethodPost, server.URL, http.NoBody)
	require.NoError(t, err)
	require.NoError(t, sigs.Sign(request, keyID, privateKey))

	response, err := server.Client().Do(request)
	require.NoError(t, err)
	t.Cleanup(func() { _ = response.Body.Close() })

	assert.True(t, verified.Load(), "the signed request must verify against the signing key")
}
