package unit

import (
	"net/http"
	"time"
)

type Time time.Time

func (value Time) IsTime() bool {
	return true
}

func (value Time) Time() time.Time {
	return time.Time(value)
}

// Get returns a value of the given property
func (value Time) Get(_ string) Value {
	return Nil{}
}

// Set returns the value with the given property set
func (value Time) Set(propertyName string, propertyValue any) Value {
	return Map{
		propertyName: propertyValue,
	}
}

// Head returns the first value in a slice
func (value Time) Head() Value {
	return value
}

// Tail returns all values in a slice except the first
func (value Time) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Time) Len() int {
	return 1
}

// IsNil returns TRUE if the value is nil
func (value Time) IsNil() bool {
	return time.Time(value).IsZero()
}

// String returs the string representation of the value
func (value Time) String() string {
	return time.Time(value).UTC().Format(http.TimeFormat)
}

// Map returns the value as a map[string]any
func (value Time) Map() map[string]any {
	return make(map[string]any)
}

// Raw returns the raw, original value
func (value Time) Raw() any {
	return time.Time(value)
}

// Clone returns a deep copy of the value
func (value Time) Clone() Value {
	return value
}
