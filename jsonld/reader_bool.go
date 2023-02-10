package jsonld

import (
	"encoding/json"
	"strconv"
	"time"
)

type Bool struct {
	value bool
}

func NewBool(value bool) Bool {
	return Bool{value: value}
}

// Property returns a sub-property of the current object
func (b Bool) Property(key string) Reader {
	return NewZero()
}

// AsBool returns the current object as a floating-point value
func (b Bool) AsBool() bool {
	return b.value
}

// AsFloat returns the current object as an integer value
func (b Bool) AsFloat() float64 {
	if b.value {
		return 1
	}
	return 0
}

// AsInt returns the current object as an integer value
func (b Bool) AsInt() int {
	if b.value {
		return 1
	}
	return 0
}

// AsString returns the current object as a string value
func (b Bool) AsString() string {
	return strconv.FormatBool(b.value)
}

// AsTime returns the current object as a time value
func (b Bool) AsTime() time.Time {
	return time.Time{}
}

// IsEmpty return TRUE if the current object is empty
func (b Bool) IsEmpty() bool {
	return b.value
}

// Tail returns a slice of all records after the first.
func (b Bool) Tail() Reader {
	return NewZero()
}

// MarshalJSON returns the JSON encoding of the current object
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}
