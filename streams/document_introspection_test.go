package streams

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDocument_Introspection confirms each Is* predicate reports TRUE for a
// document wrapping its own type and FALSE for the others. Each underlying value
// is exercised through exactly one predicate that should match.
func TestDocument_Introspection(t *testing.T) {

	// matches asserts that exactly the named predicate is true for this value.
	matches := func(name string, value any, predicate func(Document) bool) {
		t.Run(name, func(t *testing.T) {
			assert.True(t, predicate(NewDocument(value)))
		})
	}

	matches("bool", true, Document.IsBool)
	matches("int", 7, Document.IsInt)
	matches("int64", int64(7), Document.IsInt64)
	matches("float", 3.14, Document.IsFloat)
	matches("string", "hello", Document.IsString)
	matches("map", map[string]any{"a": 1}, Document.IsMap)
	matches("slice", []any{1, 2}, Document.IsSlice)
	matches("nil", nil, Document.IsNil)
}

// TestDocument_IsTyped confirms which values count as typed objects: a map with a non-empty "type" or
// "@type", and nothing else.
func TestDocument_IsTyped(t *testing.T) {

	typed := func(value any) bool {
		document := NewDocument(value)
		assert.Equal(t, !document.IsTyped(), document.NotTyped(), "NotTyped disagrees for %#v", value)
		return document.IsTyped()
	}

	assert.True(t, typed(map[string]any{"type": "Note"}))
	assert.True(t, typed(map[string]any{"type": []any{"Note", "Extra"}}))
	assert.True(t, typed(map[string]any{"@type": "Note"}))
	assert.True(t, typed(map[string]any{"type": "", "@type": "Note"}))

	assert.False(t, typed(map[string]any{"id": "https://example.com/note"}))
	assert.False(t, typed(map[string]any{"type": ""}))
	assert.False(t, typed(map[string]any{"type": []any{}}))
	assert.False(t, typed(map[string]any{"type": nil}))
	assert.False(t, typed(map[string]any{"@type": ""}))
	assert.False(t, typed(map[string]any{}))
	assert.False(t, typed("https://example.com/note"))
	assert.False(t, typed([]any{map[string]any{"type": "Note"}}))
	assert.False(t, typed(nil))
	assert.False(t, typed(42))
}

// TestDocument_IsTyped_NeverLoads confirms that asking whether a bare URL is typed never loads it,
// even when the document it points to has a type.
func TestDocument_IsTyped_NeverLoads(t *testing.T) {

	// The recording client answers this URL with a typed Note, so a load would also make it typed
	inner := newRecordingClient()
	document := NewDocument("https://example.com/note", WithClient(inner))

	assert.False(t, document.IsTyped())
	assert.True(t, document.NotTyped())
	// Any load records its URL, even one made with no options
	inner.mutex.Lock()
	defer inner.mutex.Unlock()
	assert.Empty(t, inner.options, "IsTyped loaded the URL")
}

// TestDocument_NilNotNil confirms IsNil/NotNil/NotEmpty/IsEmpty agree about
// emptiness for nil, empty, and populated documents.
func TestDocument_NilNotNil(t *testing.T) {

	nilDoc := NilDocument()
	assert.True(t, nilDoc.IsNil())
	assert.False(t, nilDoc.NotNil())
	assert.True(t, nilDoc.IsEmpty())
	assert.False(t, nilDoc.NotEmpty())

	populated := NewDocument(map[string]any{"id": "x"})
	assert.False(t, populated.IsNil())
	assert.True(t, populated.NotNil())
	assert.False(t, populated.IsEmpty())
	assert.True(t, populated.NotEmpty())
}

// TestDocument_Conversions confirms the scalar conversion accessors return the
// expected typed values.
func TestDocument_Conversions(t *testing.T) {

	t.Run("Bool", func(t *testing.T) {
		assert.True(t, NewDocument(true).Bool())
		assert.False(t, NewDocument(false).Bool())
	})

	t.Run("Int", func(t *testing.T) {
		assert.Equal(t, 42, NewDocument(42).Int())
	})

	t.Run("Float", func(t *testing.T) {
		assert.Equal(t, 3.5, NewDocument(3.5).Float())
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "hello", NewDocument("hello").String())
	})

	t.Run("Time", func(t *testing.T) {
		when := time.Date(2024, time.January, 22, 0, 0, 0, 0, time.UTC)
		// Times round-trip through the W3C/HTTP string format.
		doc := NewDocument(when.Format("Mon, 02 Jan 2006 15:04:05 GMT"))
		assert.True(t, when.Equal(doc.Time()))
	})
}

// TestDocument_Value confirms Value returns the raw unwrapped payload.
func TestDocument_Value(t *testing.T) {
	assert.Equal(t, "hello", NewDocument("hello").Value())
	assert.Equal(t, 42, NewDocument(42).Value())
	assert.Nil(t, NilDocument().Value())
}

// TestDocument_String_Sanitizes confirms String() strips HTML (bluemonday
// strict policy) while HTMLString() keeps user-generated-content-safe markup.
func TestDocument_String_Sanitizes(t *testing.T) {

	doc := NewDocument("<b>bold</b> <script>evil()</script>")

	// Strict policy removes all tags; the inner text survives.
	assert.Equal(t, "bold ", doc.String())

	// UGC policy keeps safe formatting tags but still drops scripts.
	html := doc.HTMLString()
	assert.Contains(t, html, "<b>bold</b>")
	assert.NotContains(t, html, "<script>")
}
