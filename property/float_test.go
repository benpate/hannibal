package property

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFloat exercises the full Value contract for the Float type.
func TestFloat(t *testing.T) {

	value := Float(3.14)

	assert.True(t, value.IsFloat())
	assert.Equal(t, 3.14, value.Float())
	assert.True(t, IsFloat(value))

	// Scalar Value contract.
	assert.Equal(t, Nil{}, value.Get("anything"))
	assert.Equal(t, Map{"key": "v"}, value.Set("key", "v"))
	assert.Equal(t, value, value.Head())
	assert.Equal(t, Nil{}, value.Tail())
	assert.Equal(t, 1, value.Len())
	assert.Equal(t, 3.14, value.Raw())
	assert.Equal(t, map[string]any{}, value.Map())
	assert.Equal(t, value, value.Clone())

	// String renders the float via rosetta's converter.
	assert.Equal(t, "3.14", value.String())
}

// TestFloat_IsNil confirms Float reports nil only at its zero value.
func TestFloat_IsNil(t *testing.T) {
	assert.False(t, Float(0.1).IsNil())
	assert.False(t, Float(-1).IsNil())
	assert.True(t, Float(0).IsNil())
}
