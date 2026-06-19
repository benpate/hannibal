package streams

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
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

// TestDocument_Append confirms appending two distinct values retains BOTH, in
// order, with no phantom empty element.
func TestDocument_Append(t *testing.T) {

	doc := NewDocument(mapof.NewAny())

	assert.True(t, doc.Append(vocab.PropertyTo, "https://example.com/1"))
	assert.True(t, doc.Append(vocab.PropertyTo, "https://example.com/2"))

	got := doc.Get(vocab.PropertyTo).SliceOfString()
	assert.Equal(t, []string{"https://example.com/1", "https://example.com/2"}, []string(got))
}

// TestDocument_AppendString confirms AppendString adds to a property and NOOPs
// on the empty string.
func TestDocument_AppendString(t *testing.T) {

	doc := NewDocument(mapof.NewAny())

	// Empty value is a no-op and reports false.
	assert.False(t, doc.AppendString(vocab.PropertyCC, ""))

	// Appending a real value succeeds and the value is retrievable.
	assert.True(t, doc.AppendCC("https://example.com/1"))

	assert.Equal(t, "https://example.com/1", doc.Get(vocab.PropertyCC).Head().String())
}
