package streams

import (
	"encoding/json"
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/metadata"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/require"
)

// seedJSON returns a corpus of valid, malformed, and hostile JSON payloads shared
// by the streams fuzz targets. ActivityPub documents arrive from untrusted remote
// servers, so every parser below must survive arbitrary bytes without panicking.
func seedJSON(f *testing.F) {
	f.Add([]byte(``))
	f.Add([]byte(`null`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`""`))
	f.Add([]byte(`0`))
	f.Add([]byte(`"a string"`))
	f.Add([]byte(`{"type":"Note","content":"hello","id":"https://x/1"}`))
	f.Add([]byte(`{"type":"Create","actor":"https://x/@me","object":{"type":"Note"}}`))
	f.Add([]byte(`{"type":"Collection","totalItems":2,"items":[{"id":"a"},{"id":"b"}]}`))
	f.Add([]byte(`{"type":"OrderedCollection","orderedItems":["a","b"]}`))
	f.Add([]byte(`{"@context":["https://www.w3.org/ns/activitystreams",{"toot":"x"}]}`))
	f.Add([]byte(`{"@context":"https://www.w3.org/ns/activitystreams"}`))
	f.Add([]byte(`{"to":["a","b"],"cc":"c","tag":[{"type":"Mention","href":"h"}]}`))
	f.Add([]byte(`{"items":{"id":"single-not-array"}}`))
	f.Add([]byte(`{"height":"not-a-number","width":1.5}`))
	f.Add([]byte(`{` + `"a":` + `"unterminated`))            // malformed
	f.Add([]byte(`{"deeply":{"nested":{"value":[[[1]]]}}}`)) // nesting
}

// walkDocument exercises the accessor surface of a parsed Document so that the fuzzer
// reaches code beyond the unmarshaler itself. None of these calls may panic, no matter
// what was parsed.
func walkDocument(t *testing.T, doc Document) {
	t.Helper()

	_ = doc.Type()
	_ = doc.ID()
	_ = doc.IsNil()
	_ = doc.IsString()
	_ = doc.IsMap()
	_ = doc.IsSlice()
	_ = doc.Slice()
	_ = doc.Actor().ID()
	_ = doc.Object().Type()
	_ = doc.Content()
	_ = doc.Summary()
	_ = doc.HTMLString()
	_ = doc.Height()
	_ = doc.Published()

	// Walk the addressee and item iterators, which traverse nested values.
	for range doc.RangeAddressees() {
	}
	for range doc.Range() {
	}
}

// FuzzDocumentUnmarshalJSON parses arbitrary bytes into a Document and then walks its
// accessors, confirming that neither parsing nor traversal panics on hostile input.
func FuzzDocumentUnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		document := NilDocument()

		// A parse error is an acceptable outcome; we only require no panic.
		if err := json.Unmarshal(data, &document); err != nil {
			return
		}

		walkDocument(t, document)
	})
}

// FuzzCollectionUnmarshalJSON ensures Collection.UnmarshalJSON never panics on arbitrary input.
func FuzzCollectionUnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		var collection Collection
		_ = json.Unmarshal(data, &collection)
	})
}

// FuzzOrderedCollectionUnmarshalJSON ensures OrderedCollection.UnmarshalJSON never panics.
func FuzzOrderedCollectionUnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		var collection OrderedCollection
		_ = json.Unmarshal(data, &collection)
	})
}

// FuzzCollectionPageUnmarshalJSON ensures CollectionPage.UnmarshalJSON never panics.
func FuzzCollectionPageUnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		var page CollectionPage
		_ = json.Unmarshal(data, &page)
	})
}

// FuzzOrderedCollectionPageUnmarshalJSON ensures OrderedCollectionPage.UnmarshalJSON never panics.
func FuzzOrderedCollectionPageUnmarshalJSON(f *testing.F) {

	seedJSON(f)

	f.Fuzz(func(t *testing.T, data []byte) {
		var page OrderedCollectionPage
		_ = json.Unmarshal(data, &page)
	})
}

// FuzzContextUnmarshalJSON ensures Context.UnmarshalJSON never panics. The custom decoder
// branches on the first byte (string, object, or array), so it must tolerate empty and
// malformed input safely.
func FuzzContextUnmarshalJSON(f *testing.F) {

	f.Add([]byte(``))
	f.Add([]byte(`"https://www.w3.org/ns/activitystreams"`))
	f.Add([]byte(`{"@vocab":"x","@language":"en"}`))
	f.Add([]byte(`["a",{"b":"c"}]`))
	f.Add([]byte(`[`))
	f.Add([]byte(`{`))
	f.Add([]byte(`123`))

	f.Fuzz(func(t *testing.T, data []byte) {
		var context Context
		_ = json.Unmarshal(data, &context)
	})
}

// FuzzDocumentMetadataStaysInternal confirms that no remote payload can set a document's metadata,
// and that no metadata, NoStore and Labels included, ever reaches the JSON a document marshals to.
func FuzzDocumentMetadataStaysInternal(f *testing.F) {

	seedJSON(f)

	// Hostile payloads that name the internal fields, in every spelling a decoder might honor
	f.Add([]byte(`{"id":"https://x/1","NoStore":true,"noStore":true,"nostore":true,"no_store":true}`))
	f.Add([]byte(`{"Metadata":{"NoStore":true,"Labels":[{"Value":"x","IsHidden":true}]}}`))
	f.Add([]byte(`{"metadata":{"noStore":true,"labels":[{"value":"x","isHidden":true}],"hashedId":"abc"}}`))
	f.Add([]byte(`{"HashedID":"abc","DocumentCategory":"Actor","RelationType":"Reply","RelationHref":"https://x/2","Replies":9}`))
	f.Add([]byte(`[{"NoStore":true},{"Metadata":{"NoStore":true}}]`))

	f.Fuzz(func(t *testing.T, data []byte) {

		document := NilDocument()

		if err := json.Unmarshal(data, &document); err != nil {
			return
		}

		// The parser fills only the value, never Metadata
		require.Zero(t, document.Metadata, "metadata set from %q", data)

		// Marking the document server-side changes nothing on the wire
		plain, err := json.Marshal(document)
		require.NoError(t, err)

		marked, err := json.Marshal(document.AddOptions(
			WithNoStore(),
			WithLabels(metadata.LabelSet{{Value: "Muted", IsHidden: true}}),
			WithRelation(vocab.RelationTypeReply, "https://example.com/parent"),
			WithDocumentCategory(vocab.ObjectTypeNote),
		))
		require.NoError(t, err)
		require.JSONEq(t, string(plain), string(marked))
	})
}

// FuzzOptionsClientLoad confirms that an OptionsClient passes any URL through untouched, with its
// bound options first, and that no URL a remote server supplies can change what it binds.
func FuzzOptionsClientLoad(f *testing.F) {

	f.Add("https://example.com/note")
	f.Add("https://example.com/note#main-key")
	f.Add("https://example.com/missing")
	f.Add("")
	f.Add("\x00\xff not a url")
	f.Add("https://example.com/note/" + string(make([]byte, 4096)))

	f.Fuzz(func(t *testing.T, url string) {

		inner := newRecordingClient()
		client := NewOptionsClient(inner, "bound")

		result, err := client.Load(url, "call")
		require.Equal(t, []any{"bound", "call"}, inner.optionsFor(url))

		if err != nil {
			require.True(t, derp.IsNotFound(err), "unexpected error for %q: %v", url, err)
			return
		}

		// ID() sanitizes its value as HTML, so compare the raw id the inner client returned
		require.Equal(t, client, result.Client())
		require.Equal(t, url, result.Get(vocab.PropertyID).rawString())
	})
}
