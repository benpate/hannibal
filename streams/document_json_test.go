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
// and the accessors then read it back. The document must be constructed via
// NewDocument first -- see TestDocument_UnmarshalJSON_ZeroValuePanics for the
// known bug that affects a zero-value Document.
func TestDocument_UnmarshalJSON(t *testing.T) {

	doc := NewDocument(nil)
	err := json.Unmarshal([]byte(`{"id":"https://example.com/2","name":"Alice"}`), &doc)

	require.NoError(t, err)
	assert.Equal(t, "https://example.com/2", doc.ID())
	assert.Equal(t, "Alice", doc.Name())
}

// TestDocument_UnmarshalJSON_ZeroValuePanics documents a KNOWN BUG: unmarshalling
// into a zero-value Document{} -- the idiomatic `var d Document; json.Unmarshal`
// -- panics with a nil pointer dereference, because document.value is a nil
// interface and UnmarshalJSON calls document.value.Raw() without a guard.
// See streams/document_json.go:27. The fix is a one-line nil check; until then
// callers must use NewDocument(). This test asserts the panic so the suite stays
// green while the bug is on record -- it is NOT an endorsement of the behavior.
func TestDocument_UnmarshalJSON_ZeroValuePanics(t *testing.T) {

	assert.Panics(t, func() {
		var doc Document
		_ = json.Unmarshal([]byte(`{"id":"x"}`), &doc)
	}, "KNOWN BUG: zero-value Document.UnmarshalJSON should not panic; add a nil guard in document_json.go")
}

// TestDocument_UnmarshalJSON_Invalid confirms malformed JSON returns an error
// rather than panicking.
func TestDocument_UnmarshalJSON_Invalid(t *testing.T) {

	doc := NewDocument(nil)
	err := json.Unmarshal([]byte(`{not valid json`), &doc)
	assert.Error(t, err)
}

// TestDocument_JSON_RoundTrip confirms a document survives a marshal/unmarshal
// round trip unchanged.
func TestDocument_JSON_RoundTrip(t *testing.T) {

	original := `{"id":"https://example.com/3","type":"Note","name":"Test"}`

	// Must use NewDocument; a zero-value Document panics (see ZeroValuePanics).
	doc := NewDocument(nil)
	require.NoError(t, json.Unmarshal([]byte(original), &doc))

	output, err := json.Marshal(doc)
	require.NoError(t, err)
	assert.JSONEq(t, original, string(output))
}
