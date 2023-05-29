package streams

import (
	"net/http"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/sliceof"
)

// Document represents a single ActivityStream document
// or document fragment.  Due to the flexibility of ActivityStreams
// (and JSON-LD), this may be a data structure such as a
// `map[string]any`, `[]any`, or a primitive type, like a
// `string`, `float`, `int` or `bool`.
type Document struct {
	value  any
	client Client
	header http.Header
}

// NewDocument creates a new Document object from a JSON-LD map[string]any
func NewDocument(value any, options ...Option) Document {
	result := Document{value: value}

	for _, option := range options {
		option(&result)
	}

	return result
}

// NilDocument returns a new, empty Document.
func NilDocument(options ...Option) Document {
	result := Document{}

	for _, option := range options {
		option(&result)
	}

	return result
}

/******************************************
 * Introspection Methods
 ******************************************/

func (document Document) IsString() bool {
	_, ok := document.value.(string)
	return ok
}

func (document Document) IsInt() bool {
	_, ok := document.value.(int)
	return ok
}

func (document Document) IsInt64() bool {
	_, ok := document.value.(int64)
	return ok
}

func (document Document) IsFloat() bool {
	_, ok := document.value.(float64)
	return ok
}

func (document Document) IsBool() bool {
	_, ok := document.value.(bool)
	return ok
}

func (document Document) IsNil() bool {
	return document.value == nil
}

func (document Document) NotNil() bool {
	return !document.IsNil()
}

func (document Document) IsArray() bool {
	_, ok := document.value.([]any)
	return ok
}

func (document Document) IsObject() bool {
	_, ok := document.value.(map[string]any)
	return ok
}

/******************************************
 * Getter Methods
 ******************************************/

// Value returns the generic data stored in this Document
func (document Document) Value() any {
	return document.value
}

// Get returns a sub-property of the current document
func (document Document) Get(key string) Document {

	if result := document.get(key); !result.IsNil() {
		return result
	}

	return NilDocument()
}

// Header returns the http Header associated with this document.
func (document Document) Header() http.Header {
	return document.header
}

// get does the actual work of looking up a value in
// the data structure.
func (document Document) get(key string) Document {

	switch typed := document.value.(type) {

	case string:
		if key == vocab.PropertyID {
			return document
		} else {
			object, _ := document.Load()
			return object.Get(key)
		}

	case map[string]any:
		return document.sub(typed[key])

	case mapof.Any:
		return document.sub(typed[key])

	case []any:
		if len(typed) > 0 {
			return document.sub(typed[0])
		}

	case sliceof.Any:
		if len(typed) > 0 {
			return document.sub(typed[0])
		}
	}

	return NilDocument()
}

// TODO: LOW: Add GetContext() method

/******************************************
 * Conversion Methods
 ******************************************/

func (document Document) Array() []any {

	return convert.SliceOfAny(document.value)
}

// Bool returns the current object as a floating-point value
func (document Document) Bool() bool {

	switch typed := document.value.(type) {

	case map[string]any:
		return document.Get(vocab.PropertyID).Bool()

	case []any:
		return document.Get(vocab.PropertyID).Bool()

	default:
		return convert.Bool(typed)
	}
}

// Float returns the current object as an integer value
func (document Document) Float() float64 {

	switch typed := document.value.(type) {

	case map[string]any:
		return document.Get(vocab.PropertyID).Float()

	case []any:
		return document.Get(vocab.PropertyID).Float()

	default:
		return convert.Float(typed)
	}
}

// Int returns the current object as an integer value
func (document Document) Int() int {

	switch typed := document.value.(type) {

	case map[string]any:
		return document.Get(vocab.PropertyID).Int()

	case []any:
		return document.Get(vocab.PropertyID).Int()

	default:
		return convert.Int(typed)
	}
}

// Map retrieves a JSON-LD document from a remote server, parses is, and returns a Document object.
func (document Document) Load() (Document, error) {

	const location = "hannibal.streams.Document.Map"

	switch document.value.(type) {

	case map[string]any:
		return document, nil

	case []any:
		return document.Head(), nil

	case string:
		uri := document.ID()
		return document.getClient().Load(uri)
	}

	return NilDocument(), derp.NewInternalError(location, "Document type is invalid", document.Value())
}

// String returns the current object as a string value
func (document Document) String() string {

	switch typed := document.value.(type) {

	case map[string]any:
		return document.Get(vocab.PropertyID).String()

	case []any:
		return document.Get(vocab.PropertyID).String()

	default:
		return convert.String(typed)
	}
}

// Time returns the current object as a time value
func (document Document) Time() time.Time {

	switch typed := document.value.(type) {

	case string:
		if result, err := time.Parse(time.RFC3339, typed); err == nil {
			return result
		}

	case int:
		return time.Unix(int64(typed), 0)

	case int64:
		return time.Unix(typed, 0)

	case float64:
		return time.Unix(int64(typed), 0)

	case []any:
		return document.Head().Time()

	case time.Time:
		return typed
	}

	return time.Time{}
}

/******************************************
 * List-based Methods
 ******************************************/

func (document Document) ForEach(fn func(Document)) {
	for current := document.Head(); !current.IsNil(); current = current.Tail() {
		fn(current)
	}
}

// Head returns the first object in a slice.
// For all other document types, it returns the current document.
func (document Document) Head() Document {

	if slice, ok := document.value.([]any); ok {

		if len(slice) > 0 {

			return Document{
				value:  slice[0],
				client: document.client,
			}
		}
	}

	return document
}

// Tail returns all records after the first in a slice.
// For all other document types, it returns a nil document.
func (document Document) Tail() Document {

	if slice, ok := document.value.([]any); ok {

		if len(slice) > 1 {

			return Document{
				value:  slice[1:],
				client: document.client,
			}
		}
	}

	return NilDocument()
}

// IsEmpty return TRUE if the current object is empty
func (document Document) IsEmptyTail() bool {

	if slice, ok := document.value.([]any); ok {
		return len(slice) < 2
	}

	return true
}

/******************************************
 * Helpers
 ******************************************/

func (document Document) getClient() Client {

	if document.client != nil {
		return document.client
	}

	return NewDefaultClient()
}

// sub returns a new Document with a new VALUE, all of the same OPTIONS as this original
func (document Document) sub(value any) Document {
	return Document{
		value:  value,
		client: document.client,
		header: document.header,
	}
}
