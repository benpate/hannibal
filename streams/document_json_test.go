package streams

import (
	"encoding/json"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocument_MarshalJSON confirms a document marshals to the JSON of its
// underlying value, not the wrapper struct.
func TestDocument_MarshalJSON(t *testing.T) {

	doc := NewDocument(map[string]any{vocab.PropertyID: "https://example.com/1"})

	output, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"https://example.com/1"}`, string(output))
}

// TestDocument_UnmarshalJSON confirms a document unmarshals JSON into its value
// and the accessors then read it back.
func TestDocument_UnmarshalJSON(t *testing.T) {

	doc := NewDocument(nil)
	err := json.Unmarshal([]byte(`{"id":"https://example.com/2","name":"Alice"}`), &doc)

	require.NoError(t, err)
	assert.Equal(t, "https://example.com/2", doc.ID())
	assert.Equal(t, "Alice", doc.Name())
}

// TestDocument_UnmarshalJSON_ZeroValue confirms the idiomatic
// `var d Document; json.Unmarshal(data, &d)` works without panicking. A
// zero-value Document has a nil value interface; UnmarshalJSON must guard for it.
func TestDocument_UnmarshalJSON_ZeroValue(t *testing.T) {

	var doc Document
	err := json.Unmarshal([]byte(`{"id":"https://example.com/9","name":"Bob"}`), &doc)

	require.NoError(t, err)
	assert.Equal(t, "https://example.com/9", doc.ID())
	assert.Equal(t, "Bob", doc.Name())
}

// TestDocument_UnmarshalJSON_Invalid confirms malformed JSON returns an error
// rather than panicking.
func TestDocument_UnmarshalJSON_Invalid(t *testing.T) {

	var doc Document
	err := json.Unmarshal([]byte(`{not valid json`), &doc)
	assert.Error(t, err)
}

// TestDocument_JSON_RoundTrip confirms a document survives a marshal/unmarshal
// round trip unchanged.
func TestDocument_JSON_RoundTrip(t *testing.T) {

	original := `{"id":"https://example.com/3","type":"Note","name":"Test"}`

	var doc Document
	require.NoError(t, json.Unmarshal([]byte(original), &doc))

	output, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.JSONEq(t, original, string(output))
}
