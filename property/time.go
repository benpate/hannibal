package property

import (
	"time"

	"github.com/benpate/hannibal/datetime"
)

// Time is a property Value that wraps a timestamp.
type Time time.Time

// IsTime returns TRUE if this value is a timestamp.
func (value Time) IsTime() bool {
	return true
}

// Time returns the underlying timestamp value.
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

// String returns the string representation of the value
func (value Time) String() string {

	// Timestamps serialize as AS2-conformant date-times, which returns an empty
	// string for a zero time. https://www.w3.org/TR/activitystreams-core/#dates
	return datetime.Format(time.Time(value))
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
