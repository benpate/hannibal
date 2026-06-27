package property

import "github.com/benpate/rosetta/convert"

// Int is a property Value that wraps an int.
type Int int

// IsInt returns TRUE if this value is an integer.
func (value Int) IsInt() bool {
	return true
}

// IsInt64 returns TRUE if this value is a 64-bit integer.
func (value Int) IsInt64() bool {
	return true
}

// Int returns the value as an int.
func (value Int) Int() int {
	return int(value)
}

// Int64 returns the value as an int64.
func (value Int) Int64() int64 {
	return int64(value)
}

// Get returns a value of the given property
func (value Int) Get(_ string) Value {
	return Nil{}
}

// Set returns the value with the given property set
func (value Int) Set(propertyName string, propertyValue any) Value {
	return Map{
		propertyName: propertyValue,
	}
}

// Head returns the first value in a slice
func (value Int) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value Int) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Int) Len() int {
	return 1
}

// IsNil returns TRUE if this value is nil (empty).
func (value Int) IsNil() bool {
	return value == 0
}

func (value Int) String() string {
	return convert.String(value)
}

// Map returns the value as a map[string]any.
func (value Int) Map() map[string]any {
	return make(map[string]any)
}

// Raw returns the underlying Go value.
func (value Int) Raw() any {
	return int(value)
}

// Clone returns a copy of this value.
func (value Int) Clone() Value {
	return value
}
