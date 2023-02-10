package jsonld

import (
	"encoding/json"
	"time"
)

type Map struct {
	value  map[string]any
	client *HTTPClient
}

func NewMap(value map[string]any, client *HTTPClient) Map {
	return Map{
		value:  value,
		client: client,
	}
}

// Property returns a sub-property of the current object
func (m Map) Property(key string) Reader {

	if value, ok := m.value[key]; ok {
		return NewReader(value, m.client)
	}

	return NewZero()
}

// AsBool returns the current object as a floating-point value
func (m Map) AsBool() bool {
	return false
}

// AsFloat returns the current object as an integer value
func (m Map) AsFloat() float64 {
	return 0
}

// AsInt returns the current object as an integer value
func (m Map) AsInt() int {
	return 0
}

// AsString returns the current object as a string value
func (m Map) AsString() string {
	return m.Property("id").AsString()
}

// AsTime returns the current object as a time value
func (m Map) AsTime() time.Time {
	return time.Time{}
}

// IsEmpty return TRUE if the current object is empty
func (m Map) IsEmpty() bool {
	return len(m.value) == 0
}

// Tail returns a slice of all records after the first.
func (m Map) Tail() Reader {
	return NewZero()
}

// MarshalJSON returns the JSON encoding of the current object
func (m Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.value)
}
