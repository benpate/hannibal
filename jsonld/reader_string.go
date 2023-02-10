package jsonld

import (
	"encoding/json"
	"strconv"
	"time"
)

type String struct {
	value  string
	client *HTTPClient
}

func NewString(value string, client *HTTPClient) String {
	return String{
		value:  value,
		client: client,
	}
}

// Property returns a sub-property of the current object
func (s String) Property(key string) Reader {
	return NewZero()
}

// AsBool returns the current object as a floating-point value
func (s String) AsBool() bool {
	return s.value == "true"
}

// AsFloat returns the current object as an integer value
func (s String) AsFloat() float64 {
	result, _ := strconv.ParseFloat(s.value, 64)
	return result
}

// AsInt returns the current object as an integer value
func (s String) AsInt() int {
	result, _ := strconv.Atoi(s.value)
	return result
}

// AsString returns the current object as a string value
func (s String) AsString() string {
	return s.value
}

// AsTime returns the current object as a time value
func (s String) AsTime() time.Time {
	result, _ := time.Parse(time.RFC3339, s.value)
	return result
}

// IsEmpty return TRUE if the current object is empty
func (s String) IsEmpty() bool {
	return s.value == ""
}

// Tail returns a slice of all records after the first.
func (s String) Tail() Reader {
	return NewZero()
}

// MarshalJSON returns the JSON encoding of the current object
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}
