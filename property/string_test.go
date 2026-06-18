package property

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestString exercises the full Value contract for the String type.
func TestString(t *testing.T) {

	value := NewValue("Hello, World!")

	require.Equal(t, "Hello, World!", value.Raw())
	require.Equal(t, "Hello, World!", value.Head().Raw())
	require.Nil(t, value.Tail().Raw())
	require.True(t, IsString(value))

	// String() lives on the concrete String type, not the Value interface.
	require.IsType(t, String(""), value)
	assert.Equal(t, "Hello, World!", value.(String).String())
	assert.Equal(t, 1, value.Len())
}

// TestString_Get confirms a bare string answers to its own "id" property (a
// string is treated as the id of an otherwise-empty object) and nothing else.
func TestString_Get(t *testing.T) {

	value := String("https://example.com/actor")

	// Asking for the id returns the string itself.
	assert.Equal(t, value, value.Get(vocab.PropertyID))

	// Any other property is absent.
	assert.Equal(t, Nil{}, value.Get(vocab.PropertyName))
	assert.Equal(t, Nil{}, value.Get(""))
}

// TestString_Set confirms setting a property promotes the string into a Map,
// preserving the original string as the id.
func TestString_Set(t *testing.T) {

	value := String("https://example.com/actor")
	result := value.Set(vocab.PropertyName, "Alice")

	assert.Equal(t, String("https://example.com/actor"), result.Get(vocab.PropertyID))
	assert.Equal(t, "Alice", result.Get(vocab.PropertyName).Raw())
}

// TestString_IsNil confirms String reports nil only when empty.
func TestString_IsNil(t *testing.T) {
	assert.True(t, String("").IsNil())
	assert.False(t, String("x").IsNil())
}

// TestString_Map confirms a string renders to a single-key map keyed by id.
func TestString_Map(t *testing.T) {
	value := String("hello")
	assert.Equal(t, map[string]any{vocab.PropertyID: String("hello")}, value.Map())
}

// TestString_Clone confirms a string clones to an equal value (strings are
// immutable, so the clone is trivially independent).
func TestString_Clone(t *testing.T) {
	value := String("hello")
	assert.Equal(t, value, value.Clone())
}
