package streams

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// TestNewOrderedCollection confirms the constructor sets the ID, type, default
// context, and an empty (non-nil) OrderedItems slice.
func TestNewOrderedCollection(t *testing.T) {

	collection := NewOrderedCollection("https://example.com/outbox")

	assert.Equal(t, "https://example.com/outbox", collection.ID)
	assert.Equal(t, vocab.CoreTypeOrderedCollection, collection.Type)
	assert.NotZero(t, collection.Context.Length())
	assert.NotNil(t, collection.OrderedItems)
	assert.Len(t, collection.OrderedItems, 0)
}

// TestNewOrderedCollectionPage confirms the constructor sets the ID, type,
// parent, and default context.
func TestNewOrderedCollectionPage(t *testing.T) {

	page := NewOrderedCollectionPage("https://example.com/outbox?page=1", "https://example.com/outbox")

	assert.Equal(t, "https://example.com/outbox?page=1", page.ID)
	assert.Equal(t, vocab.CoreTypeOrderedCollectionPage, page.Type)
	assert.Equal(t, "https://example.com/outbox", page.PartOf)
	assert.NotZero(t, page.Context.Length())
	assert.NotNil(t, page.OrderedItems)
}
