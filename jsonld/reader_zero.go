package jsonld

import "time"

type Zero struct{}

func NewZero() Zero {
	return Zero{}
}

// Property returns a sub-property of the current object
func (zero Zero) Property(key string) Reader {
	return zero
}

// AsBool returns the current object as a floating-point value
func (zero Zero) AsBool() bool {
	return false
}

// AsFloat returns the current object as an integer value
func (zero Zero) AsFloat() float64 {
	return 0
}

// AsInt returns the current object as an integer value
func (zero Zero) AsInt() int {
	return 0
}

// AsString returns the current object as a string value
func (zero Zero) AsString() string {
	return ""
}

// AsTime returns the current object as a time value
func (zero Zero) AsTime() time.Time {
	return time.Time{}
}

// IsEmpty return TRUE if the current object is empty
func (zero Zero) IsEmpty() bool {
	return true
}

// Tail returns a slice of all records after the first.
func (zero Zero) Tail() Reader {
	return zero
}
