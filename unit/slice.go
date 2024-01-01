package unit

type Slice []any

func (value Slice) IsSlice() bool {
	return true
}

func (value Slice) Slice() []any {
	return []any(value)
}

// Get returns a value of the given property
func (value Slice) Get(name string) Value {
	return value.Head().Get(name)
}

// Set returns the value with the given property set
func (value Slice) Set(name string, newValue any) Value {

	if len(value) == 0 {
		first := Nil{}.Set(name, newValue)
		return Slice([]any{first})
	}

	first := value.Head().Set(name, newValue)
	return Slice(append([]any{first.Raw()}, value[1:]...))
}

// Head returns the first value in a slice
func (value Slice) Head() Value {
	if len(value) == 0 {
		return Nil{}
	}

	return NewValue(value[0])
}

// Tail returns all values in a slice except the first
func (value Slice) Tail() Value {
	if len(value) == 0 {
		return Nil{}
	}

	return Slice(value[1:])
}

// Len returns the number of elements in the value
func (value Slice) Len() int {
	return len(value)
}

// IsNil returns true if the value is nil
func (value Slice) IsNil() bool {
	return len(value) == 0
}

// String returns a string representation of the value
func (value Slice) String() string {
	return "" // value.Head().String()
}

// Map returns the value as a map
func (value Slice) Map() map[string]any {
	return value.Head().Map()
}

// Raw returns the original wrapped value
func (value Slice) Raw() any {
	return []any(value)
}

// Clone returns a deep copy of the value
func (value Slice) Clone() Value {
	result := make([]any, len(value))
	copy(result, value)

	return Slice(value)
}
