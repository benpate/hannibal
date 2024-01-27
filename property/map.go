package property

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
)

type Map map[string]any

// Get returns a value of the given property
func (value Map) Get(name string) Value {

	if property, ok := value[name]; ok {
		return NewValue(property)
	}

	return Nil{}
}

// Set returns the value with the given property set
func (value Map) Set(name string, newValue any) Value {
	value[name] = newValue
	return value
}

// Head returns the first value in a slice
func (value Map) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value Map) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Map) Len() int {
	return 1
}

func (value Map) IsNil() bool {
	return len(value) == 0
}

func (value Map) String() string {
	return convert.String(value[vocab.PropertyID])
}

func (value Map) Raw() any {
	return map[string]any(value)
}

func (value Map) Clone() Value {
	result := make(map[string]any)

	for key, value := range value {
		result[key] = value
	}

	return Map(value)
}

/******************************************
 * IsMapper Interface
 ******************************************/

func (value Map) IsMap() bool {
	return true
}

func (value Map) Map() map[string]any {
	return value
}

func (value Map) MapKeys() []string {
	result := make([]string, 0, len(value))
	for key := range value {
		result = append(result, key)
	}
	return result
}
