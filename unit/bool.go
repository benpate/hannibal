package unit

import "github.com/benpate/rosetta/convert"

type Bool bool

func (value Bool) IsBool() bool {
	return true
}

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

func (value Bool) IsNil() bool {
	return false
}

func (value Bool) String() string {
	return convert.String(value)
}

func (value Bool) Map() map[string]any {
	return make(map[string]any)
}

func (value Bool) Raw() any {
	return bool(value)
}

func (value Bool) Clone() Value {
	return value
}
