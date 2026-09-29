package streams

import (
	"strings"
	"sync"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOptionsClient_PrependsOptions confirms that bound options reach the inner client ahead of the
// Load's own, in order.
func TestOptionsClient_PrependsOptions(t *testing.T) {

	inner := newRecordingClient()
	client := NewOptionsClient(inner, "bound-1", "bound-2")

	_, err := client.Load("https://example.com/note", "call-1")
	require.NoError(t, err)

	require.Equal(t, []any{"bound-1", "bound-2", "call-1"}, inner.optionsFor("https://example.com/note"))
}

// TestOptionsClient_NoBoundOptions confirms that a client with nothing bound passes a Load's options
// through unchanged.
func TestOptionsClient_NoBoundOptions(t *testing.T) {

	inner := newRecordingClient()
	client := NewOptionsClient(inner)

	_, err := client.Load("https://example.com/note", "call-1")
	require.NoError(t, err)

	require.Equal(t, []any{"call-1"}, inner.optionsFor("https://example.com/note"))
}

// TestOptionsClient_BindsResults confirms that a loaded document is bound to the OptionsClient,
// not to the client that produced it.
func TestOptionsClient_BindsResults(t *testing.T) {

	inner := newRecordingClient()
	client := NewOptionsClient(inner, "bound")

	result, err := client.Load("https://example.com/note")
	require.NoError(t, err)

	require.Equal(t, client, result.Client())
}

// TestOptionsClient_ImplicitLoadsCarryOptions confirms that a getter on a bare-URL value of a loaded
// document loads through the OptionsClient, bound options included, and so does the next level down.
func TestOptionsClient_ImplicitLoadsCarryOptions(t *testing.T) {

	inner := newRecordingClient()
	inner.documents["https://example.com/note"] = map[string]any{
		vocab.PropertyID:           "https://example.com/note",
		vocab.PropertyType:         vocab.ObjectTypeNote,
		vocab.PropertyAttributedTo: "https://example.com/alice",
	}
	inner.documents["https://example.com/alice"] = map[string]any{
		vocab.PropertyID:    "https://example.com/alice",
		vocab.PropertyType:  vocab.ActorTypePerson,
		vocab.PropertyName:  "Alice",
		vocab.PropertyImage: "https://example.com/alice/header",
	}
	inner.documents["https://example.com/alice/header"] = map[string]any{
		vocab.PropertyID:   "https://example.com/alice/header",
		vocab.PropertyType: vocab.ObjectTypeImage,
		vocab.PropertyName: "Header",
	}

	client := NewOptionsClient(inner, "bound")

	note, err := client.Load("https://example.com/note")
	require.NoError(t, err)

	// One implicit load: the author is a bare URL
	author := note.AttributedTo()
	require.Equal(t, "Alice", author.Name())
	require.Equal(t, []any{"bound"}, inner.optionsFor("https://example.com/alice"))

	// A getter on the author's own bare-URL value loads the same way
	require.Equal(t, "Header", author.Get(vocab.PropertyImage).Name())
	require.Equal(t, []any{"bound"}, inner.optionsFor("https://example.com/alice/header"))
}

// TestOptionsClient_StackedClients confirms that wrapping one OptionsClient in another passes both
// sets of options, the inner client's first, and binds results to the outer client.
func TestOptionsClient_StackedClients(t *testing.T) {

	inner := newRecordingClient()
	client := NewOptionsClient(NewOptionsClient(inner, "inner"), "outer")

	result, err := client.Load("https://example.com/note", "call")
	require.NoError(t, err)

	require.Equal(t, []any{"inner", "outer", "call"}, inner.optionsFor("https://example.com/note"))
	require.Equal(t, client, result.Client())
}

// TestOptionsClient_KeepsPrivateCopy confirms that changing the caller's slice after construction
// does not change the options this client adds.
func TestOptionsClient_KeepsPrivateCopy(t *testing.T) {

	inner := newRecordingClient()
	options := []any{"original"}
	client := NewOptionsClient(inner, options...)

	options[0] = "changed"

	_, err := client.Load("https://example.com/note")
	require.NoError(t, err)

	require.Equal(t, []any{"original"}, inner.optionsFor("https://example.com/note"))
}

// TestOptionsClient_ConcurrentLoadsDoNotShareOptions confirms that concurrent Loads with different
// options never see each other's, even when the bound slice has spare capacity.
func TestOptionsClient_ConcurrentLoadsDoNotShareOptions(t *testing.T) {

	inner := newRecordingClient()

	// Spare capacity is what an append would write into
	bound := make([]any, 1, 16)
	bound[0] = "bound"
	client := OptionsClient{inner: inner, options: bound}

	var done sync.WaitGroup

	for index := range 50 {
		done.Go(func() {
			url := "https://example.com/note/" + string(rune('A'+index))
			_, err := client.Load(url, url)
			assert.NoError(t, err)
		})
	}

	done.Wait()

	for index := range 50 {
		url := "https://example.com/note/" + string(rune('A'+index))
		require.Equal(t, []any{"bound", url}, inner.optionsFor(url))
	}
}

// TestOptionsClient_ErrorKeepsItsCode confirms that an inner failure is returned with its derp
// code intact.
func TestOptionsClient_ErrorKeepsItsCode(t *testing.T) {

	inner := newRecordingClient()
	client := NewOptionsClient(inner, "bound")

	_, err := client.Load("https://example.com/missing")

	require.Error(t, err)
	require.True(t, derp.IsNotFound(err))
}

// TestOptionsClient_NilInner confirms that a client with no inner client returns an error from
// every method rather than panicking.
func TestOptionsClient_NilInner(t *testing.T) {

	client := NewOptionsClient(nil, "bound")

	_, err := client.Load("https://example.com/note")
	require.Error(t, err)

	require.Error(t, client.Save(NewDocument(nil)))
	require.Error(t, client.Delete("https://example.com/note"))
}

// TestOptionsClient_Delegates confirms that Save, Delete, and SetRootClient all reach the inner
// client.
func TestOptionsClient_Delegates(t *testing.T) {

	inner := newRecordingClient()
	client := NewOptionsClient(inner, "bound")

	require.NoError(t, client.Save(NewDocument(map[string]any{vocab.PropertyID: "https://example.com/saved"})))
	require.Equal(t, "https://example.com/saved", inner.savedID)

	require.NoError(t, client.Delete("https://example.com/deleted"))
	require.Equal(t, "https://example.com/deleted", inner.deletedID)

	client.SetRootClient(newRecordingClient())
	require.True(t, inner.rootClientSet)
}

// TestOptionsClient_ConstructionLeavesRootAlone confirms that building an OptionsClient never
// re-roots its inner client, which may be the top of a stack already in use.
func TestOptionsClient_ConstructionLeavesRootAlone(t *testing.T) {

	inner := newRecordingClient()
	NewOptionsClient(inner, "bound")

	require.False(t, inner.rootClientSet)
}

// TestOptionsClient_NilInnerSetRootClient confirms that SetRootClient with no inner client does
// nothing rather than panicking.
func TestOptionsClient_NilInnerSetRootClient(t *testing.T) {

	client := NewOptionsClient(nil, "bound")

	require.NotPanics(t, func() {
		client.SetRootClient(newRecordingClient())
	})
}

/******************************************
 * Helpers
 ******************************************/

// recordingClient serves canned documents and records the options each URL was loaded with
type recordingClient struct {
	mutex         sync.Mutex
	documents     map[string]map[string]any
	options       map[string][]any
	savedID       string
	deletedID     string
	rootClientSet bool
}

// newRecordingClient returns a recordingClient that answers a URL under example.com/note with a
// minimal Note when it has no document for it, and 404 for everything else
func newRecordingClient() *recordingClient {
	return &recordingClient{
		documents: make(map[string]map[string]any),
		options:   make(map[string][]any),
	}
}

// Load records the options for this URL, then returns its canned document
func (client *recordingClient) Load(url string, options ...any) (Document, error) {

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.options[url] = options

	if value, found := client.documents[url]; found {
		return NewDocument(value, WithClient(client)), nil
	}

	if strings.HasPrefix(url, "https://example.com/note") {
		return NewDocument(map[string]any{vocab.PropertyID: url, vocab.PropertyType: vocab.ObjectTypeNote}, WithClient(client)), nil
	}

	return NilDocument(), derp.NotFound("recordingClient.Load", "No document", url)
}

// optionsFor returns the options the most recent Load of this URL received
func (client *recordingClient) optionsFor(url string) []any {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	return client.options[url]
}

// Save records the document's ID
func (client *recordingClient) Save(document Document) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.savedID = document.ID()
	return nil
}

// Delete records the document ID
func (client *recordingClient) Delete(documentID string) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.deletedID = documentID
	return nil
}

// SetRootClient records that it was called
func (client *recordingClient) SetRootClient(Client) {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.rootClientSet = true
}
