package property

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInt exercises the full Value contract for the Int type.
func TestInt(t *testing.T) {

	value := Int(42)

	// Int satisfies both the int and int64 introspection interfaces.
	assert.True(t, value.IsInt())
	assert.True(t, value.IsInt64())
	assert.Equal(t, 42, value.Int())
	assert.Equal(t, int64(42), value.Int64())
	assert.True(t, IsInt(value))
	assert.True(t, IsInt64(value))

	// Scalar Value contract.
	assert.Equal(t, Nil{}, value.Get("anything"))
	assert.Equal(t, Map{"key": "v"}, value.Set("key", "v"))
	assert.Equal(t, value, value.Head())
	assert.Equal(t, Nil{}, value.Tail())
	assert.Equal(t, 1, value.Len())
	assert.Equal(t, 42, value.Raw())
	assert.Equal(t, "42", value.String())
	assert.Equal(t, map[string]any{}, value.Map())
	assert.Equal(t, value, value.Clone())
}

// TestInt_IsNil confirms Int reports nil only at its zero value.
func TestInt_IsNil(t *testing.T) {
	assert.False(t, Int(1).IsNil())
	assert.False(t, Int(-1).IsNil())
	assert.True(t, Int(0).IsNil())
}
