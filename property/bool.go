package property

import "github.com/benpate/rosetta/convert"

// Bool is a property Value that wraps a boolean.
type Bool bool

// IsBool returns TRUE if this value is a boolean.
func (value Bool) IsBool() bool {
	return true
}

// Bool returns the underlying boolean value.
func (value Bool) Bool() bool {
	return bool(value)
}

// Get returns a value of the given property
func (value Bool) Get(_ string) Value {
	return Nil{}
}

// Set returns the value with the given property set
func (value Bool) Set(property string, propertyValue any) Value {
	return Map{
		property: propertyValue,
	}
}

// Head returns the first value in a slice
func (value Bool) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value Bool) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Bool) Len() int {
	return 1
}

// IsNil returns TRUE if this boolean value is nil
func (value Bool) IsNil() bool {
	return !bool(value)
}

func (value Bool) String() string {
	return convert.String(value)
}

// Map returns the value as a map[string]any.
func (value Bool) Map() map[string]any {
	return make(map[string]any)
}

// Raw returns the underlying Go value.
func (value Bool) Raw() any {
	return bool(value)
}

// Clone returns a copy of this value.
func (value Bool) Clone() Value {
	return value
}
