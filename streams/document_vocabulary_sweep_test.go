package streams

import (
	"testing"
	"time"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// TestDocument_VocabularySweep_Documents exercises the remaining thin
// Document-returning accessors not covered by the accessor sweep, confirming
// each descends to the sub-object keyed by its property name.
func TestDocument_VocabularySweep_Documents(t *testing.T) {

	id := func(name string) map[string]any {
		return map[string]any{vocab.PropertyID: "urn:" + name}
	}

	doc := NewDocument(map[string]any{
		vocab.PropertyAttributedTo:   id("attributedTo"),
		vocab.PropertyMLSKeyPackages: id("mlsKeyPackages"),
		vocab.PropertyOneOf:          id("oneOf"),
		vocab.PropertyAnyOf:          id("anyOf"),
		vocab.PropertyClosed:         id("closed"),
		vocab.PropertyPartOf:         id("partOf"),
		vocab.PropertyRel:            id("rel"),
		vocab.PropertyReplies:        id("replies"),
		vocab.PropertySubject:        id("subject"),
		vocab.PropertyPublicKey:      id("publicKey"),
	})

	check := func(name string, accessor func(Document) Document) {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, "urn:"+name, accessor(doc).ID())
		})
	}

	check("attributedTo", Document.AttributedTo)
	check("mlsKeyPackages", Document.MLSKeyPackages)
	check("oneOf", Document.OneOf)
	check("anyOf", Document.AnyOf)
	check("closed", Document.Closed)
	check("partOf", Document.PartOf)
	check("rel", Document.Rel)
	check("replies", Document.Replies)
	check("subject", Document.Subject)
	check("publicKey", Document.PublicKey)
}

// TestDocument_VocabularySweep_Strings exercises the remaining thin
// string-returning accessors.
func TestDocument_VocabularySweep_Strings(t *testing.T) {

	doc := NewDocument(map[string]any{
		vocab.PropertyMLSCiphersuite:    "ciphersuite-value",
		vocab.PropertyPublicKeyPEM:      "pem-value",
		vocab.PropertyPreferredUsername: "alice",
		vocab.PropertyEndpoints: map[string]any{
			vocab.EndpointSharedInbox: "https://example.com/inbox",
		},
	})

	assert.Equal(t, "ciphersuite-value", doc.MLSCiphersuite())
	assert.Equal(t, "pem-value", doc.PublicKeyPEM())
	assert.Equal(t, "alice", doc.Username())

	// SharedInbox reads the "sharedInbox" endpoint from within the Endpoints object.
	assert.Equal(t, "https://example.com/inbox", doc.Endpoints().SharedInbox())
}

// TestDocument_VocabularySweep_Times exercises the time-returning accessors,
// confirming an RFC3339 value is parsed and an absent property yields the zero time.
func TestDocument_VocabularySweep_Times(t *testing.T) {

	stamp := "2026-01-02T15:04:05Z"
	expected, err := time.Parse(time.RFC3339, stamp)
	assert.NoError(t, err)

	doc := NewDocument(map[string]any{
		vocab.PropertyDeleted:   stamp,
		vocab.PropertyEndTime:   stamp,
		vocab.PropertyPublished: stamp,
		vocab.PropertyStartTime: stamp,
		vocab.PropertyUpdated:   stamp,
	})

	check := func(name string, accessor func(Document) time.Time) {
		t.Run(name, func(t *testing.T) {
			assert.True(t, accessor(doc).Equal(expected))
		})
	}

	check("Deleted", Document.Deleted)
	check("EndTime", Document.EndTime)
	check("Published", Document.Published)
	check("StartTime", Document.StartTime)
	check("Updated", Document.Updated)

	// An absent time property yields the zero value.
	empty := NewDocument(map[string]any{})
	assert.True(t, empty.Updated().IsZero())
}

// TestDocument_Items confirms Items prefers "orderedItems" over "items" and
// returns a nil document when neither is present.
func TestDocument_Items(t *testing.T) {

	t.Run("prefers orderedItems", func(t *testing.T) {
		doc := NewDocument(map[string]any{
			vocab.PropertyOrderedItems: []any{"urn:ordered"},
			vocab.PropertyItems:        []any{"urn:items"},
		})
		assert.Equal(t, "urn:ordered", doc.Items().Head().ID())
	})

	t.Run("falls back to items", func(t *testing.T) {
		doc := NewDocument(map[string]any{
			vocab.PropertyItems: []any{"urn:items"},
		})
		assert.Equal(t, "urn:items", doc.Items().Head().ID())
	})

	t.Run("neither present -> nil", func(t *testing.T) {
		doc := NewDocument(map[string]any{})
		assert.True(t, doc.Items().IsNil())
	})
}
