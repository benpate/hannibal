package hannibal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/require"
)

func TestIsActivityPubContentType(t *testing.T) {
	require.True(t, IsActivityPubContentType("application/json"))
	require.True(t, IsActivityPubContentType("application/json; everything after the semicolon is ignored"))
	require.True(t, IsActivityPubContentType("application/json; whocares=notme"))
	require.True(t, IsActivityPubContentType("application/activity+json"))
	require.True(t, IsActivityPubContentType("application/activity+json; charset=utf-8"))
	require.True(t, IsActivityPubContentType("application/ld+json"))
	require.True(t, IsActivityPubContentType("application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""))

	require.False(t, IsActivityPubContentType("literally anything else"))
	require.False(t, IsActivityPubContentType("application/xml"))
	require.False(t, IsActivityPubContentType("application/xml; whocares=notme"))
	require.False(t, IsActivityPubContentType("image/webp"))
}

// TestIsActivityPubRequest confirms the request helpers read the "Accept" header.
func TestIsActivityPubRequest(t *testing.T) {

	t.Run("activitypub accept -> true", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		request.Header.Set("Accept", vocab.ContentTypeActivityPub)
		require.True(t, IsActivityPubRequest(request))
		require.False(t, NotActivityPubRequest(request))
	})

	t.Run("html accept -> false", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		request.Header.Set("Accept", "text/html")
		require.False(t, IsActivityPubRequest(request))
		require.True(t, NotActivityPubRequest(request))
	})
}

// TestIsUndoableActivity confirms undoable activity types are recognized and
// non-undoable types (e.g. Create) are not.
func TestIsUndoableActivity(t *testing.T) {

	for _, activityType := range []string{
		vocab.ActivityTypeAnnounce,
		vocab.ActivityTypeDislike,
		vocab.ActivityTypeFollow,
		vocab.ActivityTypeLike,
		vocab.ActivityTypeBlock,
	} {
		require.True(t, IsUndoableActivity(activityType), activityType)
	}

	require.False(t, IsUndoableActivity(vocab.ActivityTypeCreate))
	require.False(t, IsUndoableActivity(vocab.ActivityTypeDelete))
	require.False(t, IsUndoableActivity("SomethingMadeUp"))
}
