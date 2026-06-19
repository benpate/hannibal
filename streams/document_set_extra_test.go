package streams

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/assert"
)

// TestDocument_SetString confirms SetString writes a single string value and
// reads it back.
func TestDocument_SetString(t *testing.T) {

	doc := NewDocument(mapof.NewAny())
	ok := doc.SetString(vocab.PropertyName, "Alice")

	assert.True(t, ok)
	assert.Equal(t, "Alice", doc.Get(vocab.PropertyName).String())
}

// TestDocument_SetRecipients confirms the addressing setters (To/BTo/CC/BCC)
// each write their mapped property.
func TestDocument_SetRecipients(t *testing.T) {

	check := func(name string, set func(Document) bool, property string) {
		t.Run(name, func(t *testing.T) {
			doc := NewDocument(mapof.NewAny())
			assert.True(t, set(doc))
			assert.Equal(t, "https://example.com/r", doc.Get(property).String())
		})
	}

	check("SetTo", func(d Document) bool { return d.SetTo("https://example.com/r") }, vocab.PropertyTo)
	check("SetBTo", func(d Document) bool { return d.SetBTo("https://example.com/r") }, vocab.PropertyBTo)
	check("SetCC", func(d Document) bool { return d.SetCC("https://example.com/r") }, vocab.PropertyCC)
	check("SetBCC", func(d Document) bool { return d.SetBCC("https://example.com/r") }, vocab.PropertyBCC)
}

// TestDocument_Append_NoopOnEmpty confirms Append ignores empty input.
func TestDocument_Append_NoopOnEmpty(t *testing.T) {

	doc := NewDocument(mapof.NewAny())
	assert.False(t, doc.Append(vocab.PropertyTo, ""))
	assert.True(t, doc.Get(vocab.PropertyTo).IsNil())
}

// TestDocument_Append documents the intended contract: appending two distinct
// values should retain BOTH, with no phantom empty element.
//
// KNOWN BUG -- Append re-reads value.Head().Get(name) on each call and the prior
// append does not persist, so the first value is lost and a leading empty string
// is left behind: two appends yield ["", second] instead of [first, second].
// This test asserts the correct contract and therefore currently FAILS.
// See streams/document_set.go Append().
//
// NOTE: we inspect the RAW slice here rather than calling SliceOfString(),
// because reading the buggy ["", ...] result through the document accessors
// triggers a separate unbounded-recursion crash -- see
// TestDocument_SliceOfString_EmptyElementRecursion.
func TestDocument_Append(t *testing.T) {

	doc := NewDocument(mapof.NewAny())

	assert.True(t, doc.Append(vocab.PropertyTo, "https://example.com/1"))
	assert.True(t, doc.Append(vocab.PropertyTo, "https://example.com/2"))

	raw := convert.SliceOfString(doc.Get(vocab.PropertyTo).Value())
	assert.NotContains(t, raw, "", "Append must not leave a phantom empty element")
	assert.Contains(t, raw, "https://example.com/1", "first appended value must be retained")
	assert.Contains(t, raw, "https://example.com/2", "second appended value must be retained")
}

// TestDocument_AppendString confirms AppendString adds to a property and NOOPs
// on the empty string.
func TestDocument_AppendString(t *testing.T) {

	doc := NewDocument(mapof.NewAny())

	// Empty value is a no-op and reports false.
	assert.False(t, doc.AppendString(vocab.PropertyCC, ""))

	// Appending a real value succeeds and the value is retrievable.
	assert.True(t, doc.AppendCC("https://example.com/1"))

	// Inspect the RAW slice rather than SliceOfString(): the phantom leading ""
	// that Append leaves behind (see TestDocument_Append) makes the accessor
	// path recurse infinitely. The appended value is present in the raw data.
	raw := convert.SliceOfString(doc.Get(vocab.PropertyCC).Value())
	assert.Contains(t, raw, "https://example.com/1")
}
