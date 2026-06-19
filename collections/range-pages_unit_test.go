package collections

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// inlineCollection builds a single-page OrderedCollection with the given item IDs
// inline (no remote paging required).
func inlineCollection(itemIDs ...string) streams.Document {
	items := make([]any, 0, len(itemIDs))
	for _, id := range itemIDs {
		items = append(items, map[string]any{vocab.PropertyID: id})
	}
	return streams.NewDocument(map[string]any{
		vocab.PropertyType:  vocab.CoreTypeOrderedCollection,
		vocab.PropertyItems: items,
	})
}

// TestRangePages_SinglePage confirms a single inline collection yields exactly
// one page (itself), with no infinite looping.
func TestRangePages_SinglePage(t *testing.T) {

	collection := inlineCollection("https://example.com/1", "https://example.com/2")

	pages := 0
	for page := range RangePages(collection) {
		pages++
		assert.Equal(t, 2, page.Items().Len())
	}

	assert.Equal(t, 1, pages)
}

// TestRangePages_Empty confirms a Nil collection yields no pages.
func TestRangePages_Empty(t *testing.T) {

	pages := 0
	for range RangePages(streams.NilDocument()) {
		pages++
	}

	assert.Equal(t, 0, pages)
}

// TestRangePages_EarlyStop confirms RangePages honors an early break from the
// consumer.
func TestRangePages_EarlyStop(t *testing.T) {

	collection := inlineCollection("https://example.com/1")

	pages := 0
	for range RangePages(collection) {
		pages++
		break
	}

	assert.Equal(t, 1, pages)
}

// TestRangeDocuments_SinglePage confirms RangeDocuments yields each item in the
// page, in order.
func TestRangeDocuments_SinglePage(t *testing.T) {

	collection := inlineCollection(
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	)

	got := collectIDs(RangeDocuments(collection))
	assert.Equal(t, []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	}, got)
}

// TestRangeDocuments_Empty confirms an empty collection yields no documents.
func TestRangeDocuments_Empty(t *testing.T) {
	assert.Empty(t, collectIDs(RangeDocuments(streams.NilDocument())))
}

// TestRangeDocuments_EarlyStop confirms RangeDocuments honors an early break.
func TestRangeDocuments_EarlyStop(t *testing.T) {

	collection := inlineCollection(
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	)

	var got []string
	for document := range RangeDocuments(collection) {
		got = append(got, document.ID())
		if document.ID() == "https://example.com/2" {
			break
		}
	}

	assert.Equal(t, []string{"https://example.com/1", "https://example.com/2"}, got)
}
