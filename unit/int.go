package unit

import "github.com/benpate/rosetta/convert"

type Int int

func (value Int) IsInt() bool {
	return true
}

func (value Int) IsInt64() bool {
	return true
}

func (value Int) Int() int {
	return int(value)
}

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

func (value Int) IsNil() bool {
	return false
}

func (value Int) String() string {
	return convert.String(value)
}

func (value Int) Map() map[string]any {
	return make(map[string]any)
}

func (value Int) Raw() any {
	return int(value)
}

func (value Int) Clone() Value {
	return value
}
