package jsonld

import (
	"encoding/json"
	"strconv"
	"time"
)

type Int struct {
	value int
}

func NewInt(value int) Int {
	return Int{value: value}
}

// Get returns a sub-property of the current object
func (i Int) Get(key string) Reader {
	return NewZero()
}

// AsBool returns the current object as a floating-point value
func (i Int) AsBool() bool {
	return i.value != 0
}

// AsFloat returns the current object as an integer value
func (i Int) AsFloat() float64 {
	return float64(i.value)
}

// AsInt returns the current object as an integer value
func (i Int) AsInt() int {
	return i.value
}

// AsString returns the current object as a string value
func (i Int) AsString() string {
	return strconv.Itoa(i.value)
}

// AsTime returns the current object as a time value
func (i Int) AsTime() time.Time {
	return time.Unix(int64(i.value), 0)
}

// IsEmpty return TRUE if the current object is empty
func (i Int) IsEmpty() bool {
	return i.value == 0
}

// Tail returns a slice of all records after the first.
func (i Int) Tail() Reader {
	return NewZero()
}

// Load retrieves a remote object if the ID is available
func (i Int) Load() Reader {
	return NewZero()
}

// MarshalJSON returns the JSON encoding of the current object
func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}
