package collections

import (
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// mapClient is an in-memory streams.Client that resolves URLs from a fixed map,
// so paging tests can follow "next" links without touching the network.
type mapClient struct {
	documents map[string]map[string]any
}

func (client mapClient) SetRootClient(streams.Client) {}

func (client mapClient) Load(uri string, options ...any) (streams.Document, error) {
	if value, ok := client.documents[uri]; ok {
		return streams.NewDocument(value, streams.WithClient(client)), nil
	}
	return streams.NilDocument(), derp.Internal("collections.mapClient.Load", "Unknown URI", uri)
}

func (client mapClient) Save(streams.Document) error    { return nil }
func (client mapClient) Delete(documentID string) error { return nil }

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

// TestRangePages_InlineFirstPage reproduces the Mastodon "replies" shape: a
// Collection header whose "first" is an INLINE CollectionPage with no "id" and
// empty items, pointing via "next" to the real page that holds the items. The
// inline first page must be used as-is (not Load()ed away to Nil), so traversal
// follows "next" and reaches the items instead of stopping early.
func TestRangePages_InlineFirstPage(t *testing.T) {

	const nextURL = "https://example.com/replies?page=true"

	client := mapClient{documents: map[string]map[string]any{
		nextURL: {
			vocab.PropertyID:   nextURL,
			vocab.PropertyType: vocab.CoreTypeCollectionPage,
			vocab.PropertyItems: []any{
				map[string]any{vocab.PropertyID: "https://example.com/reply/1"},
				map[string]any{vocab.PropertyID: "https://example.com/reply/2"},
			},
		},
	}}

	collection := streams.NewDocument(map[string]any{
		vocab.PropertyID:   "https://example.com/replies",
		vocab.PropertyType: vocab.CoreTypeCollection,
		vocab.PropertyFirst: map[string]any{
			vocab.PropertyType:  vocab.CoreTypeCollectionPage,
			vocab.PropertyNext:  nextURL,
			vocab.PropertyItems: []any{}, // empty inline first page
		},
	}, streams.WithClient(client))

	got := collectIDs(RangeDocuments(collection))
	assert.Equal(t, []string{
		"https://example.com/reply/1",
		"https://example.com/reply/2",
	}, got)
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
