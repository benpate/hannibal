package jsonld

import (
	"encoding/json"
	"time"
)

type Slice struct {
	value  []any
	client *HTTPClient
}

func NewSlice(value []any, client *HTTPClient) Slice {
	return Slice{
		value:  value,
		client: client,
	}
}

// Property returns a sub-property of the current object
func (s Slice) Property(key string) Reader {
	return s.Head().Property(key)
}

// AsBool returns the current object as a floating-point value
func (s Slice) AsBool() bool {
	return s.Head().AsBool()
}

// AsFloat returns the current object as an integer value
func (s Slice) AsFloat() float64 {
	return s.Head().AsFloat()
}

// AsInt returns the current object as an integer value
func (s Slice) AsInt() int {
	return s.Head().AsInt()
}

// AsString returns the current object as a string value
func (s Slice) AsString() string {
	return s.Head().AsString()
}

// AsTime returns the current object as a time value
func (s Slice) AsTime() time.Time {
	return s.Head().AsTime()
}

// IsEmpty return TRUE if the current object is empty
func (s Slice) IsEmpty() bool {
	return len(s.value) == 0
}

// Tail returns a slice of all records after the first.
func (s Slice) Tail() Reader {

	if len(s.value) > 1 {
		return NewSlice(s.value[1:], s.client)
	}

	return NewZero()
}

// Head returns a Reader for the first record in the slice.
func (s Slice) Head() Reader {
	if len(s.value) > 0 {
		return NewReader(s.value[0], s.client)
	}

	return NewZero()
}

// MarshalJSON returns the JSON encoding of the current object
func (s Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}
