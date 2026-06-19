package streams

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// TestDocument_ID confirms ID reads the ActivityPub "id" property and returns an
// empty string when no id is present.
func TestDocument_ID(t *testing.T) {

	assert.Equal(t, "https://example.com/1",
		NewDocument(map[string]any{vocab.PropertyID: "https://example.com/1"}).ID())

	// A bare string is treated as its own id.
	assert.Equal(t, "https://example.com/2", NewDocument("https://example.com/2").ID())

	// A document with no id returns the empty string.
	assert.Equal(t, "", NewDocument(map[string]any{vocab.PropertyName: "x"}).ID())
}

// TestDocument_Type confirms Type prefers "type", falls back to "@type", and
// returns Unknown when neither is present.
func TestDocument_Type(t *testing.T) {

	assert.Equal(t, vocab.ActivityTypeCreate,
		NewDocument(map[string]any{vocab.PropertyType: vocab.ActivityTypeCreate}).Type())

	// Falls back to the JSON-LD "@type" alternate.
	assert.Equal(t, vocab.ActorTypePerson,
		NewDocument(map[string]any{vocab.PropertyType_Alternate: vocab.ActorTypePerson}).Type())

	// Neither present -> Unknown.
	assert.Equal(t, vocab.Unknown, NewDocument(map[string]any{}).Type())
}

// TestDocument_Types confirms Types returns every declared type, handling both
// single-value and slice forms, and Unknown when absent.
func TestDocument_Types(t *testing.T) {

	// Single value.
	assert.Equal(t, []string{vocab.ObjectTypeNote},
		NewDocument(map[string]any{vocab.PropertyType: vocab.ObjectTypeNote}).Types())

	// Multiple values.
	assert.Equal(t, []string{vocab.ObjectTypeNote, "Hashtag"},
		NewDocument(map[string]any{
			vocab.PropertyType: []any{vocab.ObjectTypeNote, "Hashtag"},
		}).Types())

	// Absent -> [Unknown].
	assert.Equal(t, []string{vocab.Unknown}, NewDocument(map[string]any{}).Types())
}

// TestDocument_ScalarAccessors sweeps the simple property accessors against one
// fully-populated document, confirming each reads its mapped property.
func TestDocument_ScalarAccessors(t *testing.T) {

	doc := NewDocument(map[string]any{
		vocab.PropertyName:      "Alice",
		vocab.PropertyContent:   "Hello world",
		vocab.PropertySummary:   "A summary",
		vocab.PropertyMediaType: "text/html",
		vocab.PropertyHref:      "https://example.com/href",
		vocab.PropertyURL:       "https://example.com/url",
		vocab.PropertyHeight:    100,
		vocab.PropertyWidth:     200,
		vocab.PropertyLatitude:  45.5,
		vocab.PropertyLongitude: -122.6,
	})

	assert.Equal(t, "Alice", doc.Name())
	assert.Equal(t, "Hello world", doc.Content())
	assert.Equal(t, "A summary", doc.Summary())
	assert.Equal(t, "text/html", doc.MediaType())
	assert.Equal(t, "https://example.com/href", doc.Href())
	assert.Equal(t, "https://example.com/url", doc.URL())
	assert.Equal(t, 100, doc.Height())
	assert.Equal(t, 200, doc.Width())
	assert.Equal(t, 45.5, doc.Latitude())
	assert.Equal(t, -122.6, doc.Longitude())
}

// TestDocument_URLOrID confirms URLOrID prefers the url, then falls back to id.
func TestDocument_URLOrID(t *testing.T) {

	withURL := NewDocument(map[string]any{
		vocab.PropertyID:  "https://example.com/id",
		vocab.PropertyURL: "https://example.com/url",
	})
	assert.Equal(t, "https://example.com/url", withURL.URLOrID())

	// No url -> falls back to id.
	idOnly := NewDocument(map[string]any{vocab.PropertyID: "https://example.com/id"})
	assert.Equal(t, "https://example.com/id", idOnly.URLOrID())
}

// TestDocument_DocumentAccessors confirms the Document-returning accessors
// descend into the matching sub-property.
func TestDocument_DocumentAccessors(t *testing.T) {

	doc := NewDocument(map[string]any{
		vocab.PropertyObject: map[string]any{vocab.PropertyID: "https://example.com/object"},
		vocab.PropertyTag:    map[string]any{vocab.PropertyName: "#tag"},
	})

	assert.Equal(t, "https://example.com/object", doc.Object().ID())
	assert.Equal(t, "#tag", doc.Tag().Name())

	// An absent sub-property returns a Nil document.
	assert.True(t, doc.Result().IsNil())
}
