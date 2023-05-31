package streams

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestDocument(t *testing.T) {

	d := NewDocument(map[string]any{
		"id": "https://example.com",
	})

	require.Equal(t, "https://example.com", d.ID())
}

func TestDocumentMapOfAny(t *testing.T) {

	d := NewDocument(mapof.Any{
		"id": "https://example.com",
	})

	require.True(t, d.IsMap())
	require.Equal(t, "https://example.com", d.ID())
}
