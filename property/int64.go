package property

import "github.com/benpate/rosetta/convert"

type Int64 int64

func (value Int64) IsInt() bool {
	return true
}

func (value Int64) IsInt64() bool {
	return true
}

func (value Int64) Int() int {
	return int(value)
}

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

func (value Int64) IsNil() bool {
	return value == 0
}

func (value Int64) String() string {
	return convert.String(value)
}

func (value Int64) Map() map[string]any {
	return make(map[string]any)
}

func (value Int64) Raw() any {
	return int64(value)
}

func (value Int64) Clone() Value {
	return value
}
