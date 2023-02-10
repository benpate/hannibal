package jsonld

import (
	"encoding/json"
	"strconv"
	"time"
)

type Float struct {
	value float64
}

func NewFloat(value float64) Float {
	return Float{value: value}
}

// Property returns a sub-property of the current object
func (f Float) Property(key string) Reader {
	return NewZero()
}

// AsBool returns the current object as a floating-point value
func (f Float) AsBool() bool {
	return f.value != 0
}

// AsFloat returns the current object as an integer value
func (f Float) AsFloat() float64 {
	return f.value
}

// AsInt returns the current object as an integer value
func (f Float) AsInt() int {
	return int(f.value)
}

// AsString returns the current object as a string value
func (f Float) AsString() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

// AsTime returns the current object as a time value
func (f Float) AsTime() time.Time {
	return time.Unix(int64(f.value), 0)
}

// IsEmpty return TRUE if the current object is empty
func (f Float) IsEmpty() bool {
	return f.value == 0
}

// Tail returns a slice of all records after the first.
func (f Float) Tail() Reader {
	return NewZero()
}

// MarshalJSON returns the JSON encoding of the current object
func (f Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}
