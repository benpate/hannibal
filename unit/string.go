package unit

import "github.com/benpate/hannibal/vocab"

type String string

func (value String) IsString() bool {
	return true
}

// Get returns a value of the given property
func (value String) Get(propertyName string) Value {

	if propertyName == vocab.PropertyID {
		return value
	}

	return Nil{}
}

// Set returns the value with the given property set
func (value String) Set(propertyName string, propertyValue any) Value {
	result := Map{
		vocab.PropertyID: value,
	}

	return result.Set(propertyName, propertyValue)
}

// Head returns the first value in a slice
func (value String) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value String) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value String) Len() int {
	return 1
}

// IsNil returns TRUE if the value is nil
func (value String) IsNil() bool {
	return value == ""
}

// String returns a string representation of the value
func (value String) String() string {
	return string(value)
}

// Map returns the value as a map
func (value String) Map() map[string]any {
	return map[string]any{
		vocab.PropertyID: value,
	}
}

// Raw returns the raw, original value
func (value String) Raw() any {
	return string(value)
}

// Clone returns a deep copy of the value
func (value String) Clone() Value {
	return value
}
