package unit

import "time"

/******************************************
 * Introspection Interfaces
 ******************************************/

type IsBooler interface {
	IsBool() bool
	Bool() bool
}

type IsInter interface {
	IsInt() bool
	Int() int
}

type IsInt64er interface {
	IsInt64() bool
	Int64() int64
}

type IsFloater interface {
	IsFloat() bool
	Float() float64
}

type IsMapper interface {
	IsMap() bool
	Map() map[string]any
}

type IsSlicer interface {
	IsSlice() bool
	Slice() []any
}

type IsStringer interface {
	IsString() bool
	String() string
}

type IsTimeer interface {
	IsTime() bool
	Time() time.Time
}

/******************************************
 * Getter Interfaces
 ******************************************/

// BoolGetter is an optional interface that should be implemented
// by any unit.Value that contains a bool
type BoolGetter interface {
	// Bool returns a value typed as a bool
	Bool() bool
}

// IntGetter is an optional interface that should be implemented
// by any unit.Value that contains an int
type IntGetter interface {
	// Int returns the value typed as an int
	Int() int
}

// Int64Getter is an optional interface that should be implemented
// by any unit.Value that contains an int64
type Int64Getter interface {
	// Int64 returns the value typed as an int64
	Int64() int64
}

// FloatGetter is an optional interface that should be implemented
// by any unit.Value that contains a float64
type FloatGetter interface {
	// Float returns the value typed as a float64
	Float() float64
}

// MapGetter is an optional interface that should be implemented
// by any unit.Value that contains a map[string]any
type MapGetter interface {
	// Map returns the value typed as a map[string]any
	Map() map[string]any
}

// SliceGetter is an optional interface that should be implemented
// by any unit.Value that contains a []any
type SliceGetter interface {
	// Slice returns the value typed as a []any
	Slice() []any
}

// StringGetter is an optional interface that should be implemented
// by any unit.Value that contains a string
type StringGetter interface {
	// String returns the value typed as a string
	String() string
}

// TimeGetter is an optional interface that should be implemented
// by any unit.Value that contains a time.Time
type TimeGetter interface {
	// Time returns the value typed as a time.Time
	Time() time.Time
}
