package streams

import (
	"bytes"
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
)

type Context []ContextEntry

func NewContext(args ...string) Context {
	result := make(Context, len(args))

	for index, arg := range args {
		result[index] = NewContextEntry(arg)
	}

	return result
}

// DefaultContext represents the standard context defined by the W3C
func DefaultContext() Context {
	return NewContext(vocab.NamespaceActivityStreams)
}

func (c Context) Length() int {
	if c == nil {
		return 0
	}

	return len(c)
}

func (c Context) IsEmpty() bool {
	return c.Length() == 0
}

func (c Context) IsEmptyTail() bool {
	return c.Length() <= 1
}

func (c Context) Head() *ContextEntry {
	if c.Length() == 0 {
		return nil
	}

	return &(c[0])
}

func (c Context) Tail() Context {
	if c.Length() == 0 {
		return c
	}

	return c[1:]
}

// Add puts a new ContextEntry into the list and
// returns a pointer to it so that additional properties
// can be set.
func (c *Context) Add(vocabulary string) *ContextEntry {
	entry := NewContextEntry(vocabulary)
	*c = append(*c, entry)
	return &((*c)[len(*c)-1])
}

func (c Context) MarshalJSON() ([]byte, error) {

	const location = "writer.Context.MarshalJSON"

	switch len(c) {

	case 0:
		return []byte("null"), nil

	case 1:
		return json.Marshal(c[0])
	}

	// Otherwise, write the Context as an array
	var buffer bytes.Buffer

	buffer.WriteByte('[')

	for index, context := range c {
		if index > 0 {
			buffer.WriteByte(',')
		}

		item, err := json.Marshal(context)

		if err != nil {
			return nil, derp.Wrap(err, location, "Failed to marshal context")
		}

		buffer.Write(item)
	}

	buffer.WriteByte(']')

	return buffer.Bytes(), nil
}

func (c *Context) UnmarshalJSON(data []byte) error {

	const location = "writer.Context.UnmarshalJSON"

	// If the data is empty, then this object is empty, too
	if len(data) == 0 {
		*c = make(Context, 0)
		return nil
	}

	// If this looks like a single item, then unmarshal it as a single item
	if (data[0] == '{') || (data[0] == '"') {

		onlyContext := ContextEntry{}

		if err := json.Unmarshal(data, &onlyContext); err != nil {
			return derp.Wrap(err, location, "Failed to unmarshal context")
		}

		*c = Context{onlyContext}
		return nil
	}

	// Otherwise, this looks like an array of contexts
	var entries []ContextEntry

	if err := json.Unmarshal(data, &entries); err != nil {
		return derp.Wrap(err, location, "Failed to unmarshal context array")
	}

	*c = entries
	return nil
}
