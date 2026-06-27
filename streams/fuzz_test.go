package streams

import (
	"encoding/json"
	"testing"
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
