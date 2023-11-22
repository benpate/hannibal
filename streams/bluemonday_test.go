package streams

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBlueMonday tests the content filtering implementation in streams.Documents
func TestBlueMonday(t *testing.T) {

	badValue := map[string]any{
		"name":    "John <i>Connor</i>",
		"summary": "This is a <b>bad</b> summary <script>alert('hey there')</script>",
		"content": "<p>Some of this content should be <b>visible</b>.</p><script>alert('but not this')</script>",
	}

	badDocument := NewDocument(badValue)

	require.Equal(t, "John Connor", badDocument.Name())
	require.Equal(t, "This is a bad summary ", badDocument.Summary())
	require.Equal(t, "<p>Some of this content should be <b>visible</b>.</p>", badDocument.Content())
}
