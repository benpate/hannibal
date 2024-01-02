package property

import (
	"time"

	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Value is a wrapper for any kind of value that might be used in a streams.Document
type Value interface {

	// Get returns a value of the given property
	Get(string) Value

	// Set returns the value with the given property set
	Set(string, any) Value

	// Head returns the first value in a slices, or the value itself if it is not a slice
	Head() Value

	// Tail returns all values in a slice except the first
	Tail() Value

	// Len returns the number of elements in the value
	Len() int

	// IsNil returns TRUE if the value is empty
	IsNil() bool

	// Map returns the map representation of this value
	Map() map[string]any

	// Raw returns the raw, unwrapped value being stored
	Raw() any

	// Clone returns a deep copy of a value
	Clone() Value
}

func NewValue(value any) Value {

	switch typed := value.(type) {

	case Value:
		return typed

	case bool:
		return Bool(typed)

	case float32:
		return Float(typed)

	case float64:
		return Float(typed)

	case int:
		return Int(typed)

	case int64:
		return Int64(typed)

	case map[string]any:
		return Map(typed)

	case mapof.Any:
		return Map(typed)

	case []any:
		return Slice(typed)

	case sliceof.Any:
		return Slice(typed)

	case string:
		return String(typed)

	case time.Time:
		return Time(typed)
	}

	// More checks for wayward values (like primitive.A)

	if convert.IsMap(value) {
		return Map(convert.MapOfAny(value))
	}

	if convert.IsSlice(value) {
		return Slice(convert.SliceOfAny(value))
	}

	return Nil{}
}

/****************************************************
 * Introspection Functions
 ****************************************************/

// IsBool returns TRUE if the value represents a bool
func IsBool(value any) bool {
	if is, ok := value.(IsBooler); ok {
		return is.IsBool()
	}
	return false
}

// IsInt returns TRUE if the value represents a float
func IsFloat(value any) bool {
	if is, ok := value.(IsFloater); ok {
		return is.IsFloat()
	}
	return false
}

// IsInt returns TRUE if the value represents an int
func IsInt(value any) bool {
	if is, ok := value.(IsInter); ok {
		return is.IsInt()
	}
	return false
}

// IsInt64 returns TRUE if the value represents an int64
func IsInt64(value any) bool {
	if is, ok := value.(IsInt64er); ok {
		return is.IsInt64()
	}
	return false
}

// IsMap returns TRUE if the value represents a map
func IsMap(value any) bool {
	if is, ok := value.(IsMapper); ok {
		return is.IsMap()
	}
	return false
}

// IsSlice returns TRUE if the value represents a slice
func IsSlice(value any) bool {
	if is, ok := value.(IsSlicer); ok {
		return is.IsSlice()
	}
	return false
}

// IsString returns TRUE if the value represents a string
func IsString(value any) bool {
	if is, ok := value.(IsStringer); ok {
		return is.IsString()
	}
	return false
}

// IsTime returns TRUE if the value represents a time.Time
func IsTime(value any) bool {
	if is, ok := value.(IsTimeer); ok {
		return is.IsTime()
	}
	return false
}
