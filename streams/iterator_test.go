package streams

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollection(t *testing.T) {

	document := NewDocument("https://demo/collection", WithClient(newTestClient()))

	iterator, err := NewIterator(document)

	require.Nil(t, err)
	require.Equal(t, 3, iterator.TotalItems())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/collection-url-1", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/collection-url-2", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/collection-url-3", iterator.Next().Value())

	require.False(t, iterator.HasNext())
}

func TestOrderedCollection(t *testing.T) {

	document := NewDocument("https://demo/ordered", WithClient(newTestClient()))

	iterator, err := NewIterator(document)

	require.Nil(t, err)
	require.Equal(t, 3, iterator.TotalItems())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-1", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-2", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-3", iterator.Next().Value())

	require.False(t, iterator.HasNext())
}

func TestOrderedCollectionPaging(t *testing.T) {

	document := NewDocument("https://demo/ordered-page", WithClient(newTestClient()))

	iterator, err := NewIterator(document)

	require.Nil(t, err)
	require.Equal(t, 9, iterator.TotalItems())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-1", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-2", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-3", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-4", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-5", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-6", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-7", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-8", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-9", iterator.Next().Value())

	require.False(t, iterator.HasNext())
}

// TestOrderedColleciton_Interminus tests an (improperly?)  terminated ordered collection
// whose last page is empty.  It could happen..
func TestOrderedCollection_Interminus(t *testing.T) {

	document := NewDocument("https://demo/interminus", WithClient(newTestClient()))

	iterator, err := NewIterator(document)

	require.Nil(t, err)
	require.Equal(t, 9, iterator.TotalItems())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-1", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-2", iterator.Next().Value())

	require.True(t, iterator.HasNext())
	require.Equal(t, "https://example/url-3", iterator.Next().Value())

	require.False(t, iterator.HasNext())
}

/******************************************
 * Test API calls WITHOUT using .HasNext()
 ******************************************/

func TestCollection_BadAPI(t *testing.T) {

	document := NewDocument("https://demo/collection", WithClient(newTestClient()))

	iterator, err := NewIterator(document)

	require.Nil(t, err)

	require.Equal(t, "https://example/collection-url-1", iterator.Next().Value())
	require.Nil(t, iterator.Error())

	require.Equal(t, "https://example/collection-url-2", iterator.Next().Value())
	require.Nil(t, iterator.Error())

	require.Equal(t, "https://example/collection-url-3", iterator.Next().Value())
	require.Nil(t, iterator.Error())

	require.False(t, iterator.HasNext())
	require.Error(t, iterator.Error())
}

func TestOrderedCollectionPaging_BadAPI(t *testing.T) {

	document := NewDocument("https://demo/ordered-page", WithClient(newTestClient()))

	iterator, err := NewIterator(document)

	require.Nil(t, err)
	require.Equal(t, "https://example/url-1", iterator.Next().Value())
	require.Equal(t, "https://example/url-2", iterator.Next().Value())
	require.Equal(t, "https://example/url-3", iterator.Next().Value())
	require.Equal(t, "https://example/url-4", iterator.Next().Value())
	require.Equal(t, "https://example/url-5", iterator.Next().Value())
	require.Equal(t, "https://example/url-6", iterator.Next().Value())
	require.Equal(t, "https://example/url-7", iterator.Next().Value())
	require.Equal(t, "https://example/url-8", iterator.Next().Value())
	require.Equal(t, "https://example/url-9", iterator.Next().Value())
	require.Equal(t, "", iterator.Next().String())

	require.False(t, iterator.HasNext())
}

/******************************************
 * Test Other Errors
 ******************************************/

func TestCollectionError(t *testing.T) {

	document := NewDocument("https://missing-document/collection", WithClient(newTestClient()))

	_, err := NewIterator(document)

	require.Error(t, err)
}
