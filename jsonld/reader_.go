package jsonld

import "time"

type Reader interface {

	// Property returns a sub-property of the current object
	Property(key string) Reader

	// AsBool returns the current object as a floating-point value
	AsBool() bool

	// AsFloat returns the current object as an integer value
	AsFloat() float64

	// AsInt returns the current object as an integer value
	AsInt() int

	// AsString returns the current object as a string value
	AsString() string

	// AsTime returns the current object as a time value
	AsTime() time.Time

	// IsEmpty return TRUE if the current object is empty
	IsEmpty() bool

	// Tail returns a slice of all records after the first.
	Tail() Reader
}
