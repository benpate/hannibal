package property

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBool exercises the full Value contract for the Bool type.
func TestBool(t *testing.T) {

	value := Bool(true)

	// Type introspection.
	assert.True(t, value.IsBool())
	assert.True(t, value.Bool())
	assert.True(t, IsBool(value))

	// Scalar values have no addressable sub-properties.
	assert.Equal(t, Nil{}, value.Get("anything"))

	// Set on a scalar promotes it to a Map keyed by the property name.
	assert.Equal(t, Map{"key": "newValue"}, value.Set("key", "newValue"))

	// Head/Tail/Len: a scalar is a single element.
	assert.Equal(t, value, value.Head())
	assert.Equal(t, Nil{}, value.Tail())
	assert.Equal(t, 1, value.Len())

	// Raw round-trips to the underlying bool.
	assert.Equal(t, true, value.Raw())
	assert.Equal(t, "true", value.String())
	assert.Equal(t, map[string]any{}, value.Map())
	assert.Equal(t, value, value.Clone())
}

// TestBool_IsNil confirms Bool reports nil only when false (its zero value).
func TestBool_IsNil(t *testing.T) {
	assert.False(t, Bool(true).IsNil())
	assert.True(t, Bool(false).IsNil())
}
