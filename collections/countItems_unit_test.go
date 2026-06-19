package collections

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCountItems_TotalItems confirms the fast path: when the collection reports
// a positive TotalItems, that value is returned without paging.
func TestCountItems_TotalItems(t *testing.T) {

	collection := streams.NewDocument(map[string]any{
		vocab.PropertyType:       vocab.CoreTypeOrderedCollection,
		vocab.PropertyTotalItems: 42,
	})

	count, err := CountItems(collection)
	require.NoError(t, err)
	assert.Equal(t, 42, count)
}

// TestCountItems_Nil confirms a Nil collection counts as zero items.
func TestCountItems_Nil(t *testing.T) {

	count, err := CountItems(streams.NilDocument())
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

// TestCountItems_ByPaging confirms that when no TotalItems is reported, the items
// are counted by walking the (single, inline) page.
func TestCountItems_ByPaging(t *testing.T) {

	collection := inlineCollection(
		"https://example.com/1",
		"https://example.com/2",
		"https://example.com/3",
	)

	count, err := CountItems(collection)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

// TestCountItems_EmptyPage confirms a collection with no items and no TotalItems
// counts as zero.
func TestCountItems_EmptyPage(t *testing.T) {

	collection := streams.NewDocument(map[string]any{
		vocab.PropertyType:  vocab.CoreTypeOrderedCollection,
		vocab.PropertyItems: []any{},
	})

	count, err := CountItems(collection)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
