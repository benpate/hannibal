package property

type Nil struct{}

// Get returns a value of the given property
func (value Nil) Get(string) Value {
	return Nil{}
}

// Set returns the value with the given property set
func (value Nil) Set(string, any) Value {
	return Nil{}
}

// Head returns the first value in a slice
func (value Nil) Head() Value {
	return Nil{}
}

// Tail returns all values in a slice except the first
func (value Nil) Tail() Value {
	return Nil{}
}

// Len returns the number of elements in the value
func (value Nil) Len() int {
	return 0
}

func (value Nil) IsNil() bool {
	return true
}

func (value Nil) String() string {
	return ""
}

func (value Nil) Map() map[string]any {
	return make(map[string]any)
}

func (value Nil) Raw() any {
	return nil
}

func (value Nil) Clone() Value {
	return Nil{}
}
