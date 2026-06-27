package property

import "github.com/benpate/rosetta/convert"

// Int64 is a property Value that wraps an int64.
type Int64 int64

// IsInt returns TRUE if this value is an integer.
func (value Int64) IsInt() bool {
	return true
}

// IsInt64 returns TRUE if this value is a 64-bit integer.
func (value Int64) IsInt64() bool {
	return true
}

// Int returns the value as an int.
func (value Int64) Int() int {
	return int(value)
}

// Int64 returns the value as an int64.
func (value Int64) Int64() int64 {
	return int64(value)
}

// Get returns a value of the given property
func (value Int64) Get(_ string) Value {
	return Nil{}
}

// Set returns the value with the given property set
func (value Int64) Set(propertyName string, propertyValue any) Value {
	return Map{
		propertyName: propertyValue,
	}
}

// Head returns the first value in a slice
func (value Int64) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value Int64) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Int64) Len() int {
	return 1
}

// IsNil returns TRUE if this value is nil (empty).
func (value Int64) IsNil() bool {
	return value == 0
}

func (value Int64) String() string {
	return convert.String(value)
}

// Map returns the value as a map[string]any.
func (value Int64) Map() map[string]any {
	return make(map[string]any)
}

// Raw returns the underlying Go value.
func (value Int64) Raw() any {
	return int64(value)
}

// Clone returns a copy of this value.
func (value Int64) Clone() Value {
	return value
}
