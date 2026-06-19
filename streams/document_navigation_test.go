package streams

import (
	"sort"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
	"github.com/stretchr/testify/assert"
)

// TestDocument_Get confirms Get descends into map properties and returns a Nil
// document for absent keys.
func TestDocument_Get(t *testing.T) {

	doc := NewDocument(map[string]any{
		vocab.PropertyID:   "https://example.com/1",
		vocab.PropertyName: "Alice",
	})

	assert.Equal(t, "Alice", doc.Get(vocab.PropertyName).String())
	assert.Equal(t, "https://example.com/1", doc.Get(vocab.PropertyID).String())

	// An absent key yields a Nil document, not a panic.
	assert.True(t, doc.Get(vocab.PropertySummary).IsNil())
}

// TestDocument_Get_StringID confirms that asking a bare-string document for its
// id returns the string itself, without trying to load anything remotely.
func TestDocument_Get_StringID(t *testing.T) {

	doc := NewDocument("https://example.com/actor")
	assert.Equal(t, "https://example.com/actor", doc.Get(vocab.PropertyID).String())
}

// TestDocument_HeadTail confirms Head/Tail/Len semantics across a slice.
func TestDocument_HeadTail(t *testing.T) {

	doc := NewDocument([]any{"first", "second", "third"})

	assert.Equal(t, 3, doc.Len())
	assert.Equal(t, "first", doc.Head().String())

	tail := doc.Tail()
	assert.Equal(t, 2, tail.Len())
	assert.Equal(t, "second", tail.Head().String())

	// IsEmptyTail is true once fewer than two items remain.
	assert.False(t, doc.IsEmptyTail())
	assert.True(t, doc.Head().IsEmptyTail())
}

// TestDocument_HeadTail_Scalar confirms a scalar behaves like a single-element
// list: Head is itself, Tail is empty, Len is 1.
func TestDocument_HeadTail_Scalar(t *testing.T) {

	doc := NewDocument("only")

	assert.Equal(t, 1, doc.Len())
	assert.Equal(t, "only", doc.Head().String())
	assert.True(t, doc.Tail().IsNil())
	assert.True(t, doc.IsEmptyTail())
}

// TestDocument_Slice confirms Slice returns the raw slice and SliceOfDocuments /
// SliceOfString project it into the respective forms.
func TestDocument_Slice(t *testing.T) {

	doc := NewDocument([]any{
		map[string]any{vocab.PropertyID: "https://example.com/1"},
		map[string]any{vocab.PropertyID: "https://example.com/2"},
	})

	assert.Len(t, doc.Slice(), 2)

	docs := doc.SliceOfDocuments()
	assert.Len(t, docs, 2)
	assert.Equal(t, "https://example.com/1", docs[0].ID())

	ids := doc.SliceOfString()
	assert.Equal(t, []string{"https://example.com/1", "https://example.com/2"}, []string(ids))
}

// TestDocument_Map confirms Map renders a document to a map and honors the strip
// options.
func TestDocument_Map(t *testing.T) {

	doc := NewDocument(map[string]any{
		vocab.AtContext:    "https://www.w3.org/ns/activitystreams",
		vocab.PropertyID:   "https://example.com/1",
		vocab.PropertyTo:   "https://example.com/followers",
		vocab.PropertyName: "Alice",
	})

	// Without options, all keys survive.
	full := doc.Map()
	assert.Equal(t, "Alice", full[vocab.PropertyName])
	assert.Contains(t, full, vocab.AtContext)

	// StripContext removes the JSON-LD @context.
	stripped := doc.Map(OptionStripContext)
	assert.NotContains(t, stripped, vocab.AtContext)

	// StripRecipients removes the addressing fields.
	noRecipients := doc.Map(OptionStripRecipients)
	assert.NotContains(t, noRecipients, vocab.PropertyTo)
}

// TestDocument_Map_String confirms a bare-string document maps to a single id key.
func TestDocument_Map_String(t *testing.T) {

	doc := NewDocument("https://example.com/actor")
	result := doc.Map()
	// A bare string maps to its id; the stored value is a property.String, so
	// compare via fmt rather than asserting the concrete wrapper type.
	assert.Equal(t, "https://example.com/actor", convert.String(result[vocab.PropertyID]))
}

// TestDocument_MapKeys returns all keys of a map document (order-independent).
func TestDocument_MapKeys(t *testing.T) {

	doc := NewDocument(map[string]any{"a": 1, "b": 2, "c": 3})
	keys := doc.MapKeys()

	sort.Strings(keys)
	assert.Equal(t, []string{"a", "b", "c"}, []string(keys))

	// A non-map document has no keys.
	assert.Empty(t, NewDocument("string").MapKeys())
}

// TestDocument_SliceOfString_EmptyElementRecursion documents a KNOWN BUG that
// is intentionally NOT executed: a slice document containing an empty-string
// element overflows the stack when projected with SliceOfString (or any path
// that calls ID() on the empty element).
//
// Minimal reproducer (DO NOT run -- it is a fatal, unrecoverable stack overflow,
// not a recoverable panic):
//
//	NewDocument([]any{"", "https://example.com/1"}).SliceOfString()
//
// Root cause: Document.Get on an empty string enters the remote Load() path
// (document_.go:121); Load() calls ID(), ID() calls Get(), and the cycle never
// terminates because the empty id never short-circuits in the slice context.
// The empty element commonly arises from the Append bug (see TestDocument_Append).
// This test is skipped so it records the bug without crashing the suite.
func TestDocument_SliceOfString_EmptyElementRecursion(t *testing.T) {
	t.Skip("KNOWN BUG: SliceOfString on a slice with an empty-string element " +
		"recurses infinitely (fatal stack overflow). See streams/document_.go Get/Load/ID. " +
		"Not executed because the overflow is unrecoverable and would crash the test binary.")
}

// TestDocument_Clone confirms a cloned document carries an independent value:
// mutating the clone must not change the original.
func TestDocument_Clone(t *testing.T) {

	original := NewDocument(map[string]any{vocab.PropertyName: "original"})
	clone := original.Clone()

	clone.SetProperty(vocab.PropertyName, "MUTATED")

	assert.Equal(t, "original", original.Get(vocab.PropertyName).String())
}
