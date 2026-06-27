package property

import "github.com/benpate/rosetta/convert"

// Float is a property Value that wraps a float64.
type Float float64

// IsFloat returns TRUE if this value is a floating-point number.
func (value Float) IsFloat() bool {
	return true
}

// Float returns the underlying floating-point value.
func (value Float) Float() float64 {
	return float64(value)
}

// Get returns a value of the given property
func (value Float) Get(_ string) Value {
	return Nil{}
}

// Set returns the value with the given property set
func (value Float) Set(propertyName string, propertyValue any) Value {
	return Map{
		propertyName: propertyValue,
	}
}

// Head returns the first value in a slice
func (value Float) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value Float) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Float) Len() int {
	return 1
}

// IsNil returns TRUE if this value is nil (empty).
func (value Float) IsNil() bool {
	return value == 0
}

func (value Float) String() string {
	return convert.String(value)
}

// Map returns the value as a map[string]any.
func (value Float) Map() map[string]any {
	return make(map[string]any)
}

// Raw returns the underlying Go value.
func (value Float) Raw() any {
	return float64(value)
}

// Clone returns a copy of this value.
func (value Float) Clone() Value {
	return value
}
