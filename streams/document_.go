package streams

import (
	"strconv"
	"time"
)

// Document represents a single ActivityStream document
// or document fragment.  Due to the flexibility of ActivityStreams
// (and JSON-LD), this may be a data structure such as a
// `map[string]any`, `[]any`, or a primitive type, like a
// `string`, `float`, `int` or `bool`.
type Document struct {
	value any
	cache Cache
}

// NewDocument creates a new Document object from a JSON-LD map[string]any
func NewDocument(value map[string]any, cache Cache) Document {
	return Document{
		value: value,
		cache: cache,
	}
}

func NewID(value string, cache Cache) Document {
	return Document{
		value: value,
		cache: cache,
	}
}

// NilDocument returns a new, empty Document.
func NilDocument() Document {
	return Document{}
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

	// Look for the value in the document.  This should
	// happen 99.9% of the time.
	if result := document.get(key); !result.IsNil() {
		return result
	}

	// Odd case: Search alternate values that ActivityPub
	// has aliased with/without "@" signs.  Grr..
	switch key {
	case "@id":
		return document.get("id")
	case "@type":
		return document.get("type")
	case "id":
		return document.get("@id")
	case "type":
		return document.get("@type")
	default:
		return NilDocument()
	}
}

// get does the actual work of looking up a value in
// the data structure.
func (document Document) get(key string) Document {

	switch typed := document.value.(type) {

	case string:
		if key == "id" {
			return document
		} else {
			object, _ := document.AsObject()
			return object.Get(key)
		}

	case map[string]any:
		return document.sub(typed[key])

	case []any:
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

// AsBool returns the current object as a floating-point value
func (document Document) AsBool() bool {

	switch typed := document.value.(type) {

	case string:
		return typed == "true"

	case int:
		return typed != 0

	case int64:
		return typed != 0

	case float64:
		return typed != 0

	case bool:
		return typed

	case map[string]any:
		return document.Get("id").AsBool()

	case []any:
		return document.Get("id").AsBool()
	}

	return false
}

// AsFloat returns the current object as an integer value
func (document Document) AsFloat() float64 {

	switch typed := document.value.(type) {

	case string:
		if result, err := strconv.ParseFloat(typed, 64); err != nil {
			return result
		}

	case int:
		return float64(typed)

	case int64:
		return float64(typed)

	case float64:
		return typed

	case bool:
		if typed {
			return 1
		}
		return 0

	case map[string]any:
		return document.Get("id").AsFloat()

	case []any:
		return document.Get("id").AsFloat()
	}

	return 0
}

// AsInt returns the current object as an integer value
func (document Document) AsInt() int {

	switch typed := document.value.(type) {

	case string:
		if result, err := strconv.Atoi(typed); err != nil {
			return result
		}

	case int:
		return typed

	case int64:
		return int(typed)

	case float64:
		return int(typed)

	case bool:
		if typed {
			return 1
		}
		return 0

	case map[string]any:
		return document.Get("id").AsInt()

	case []any:
		return document.Get("id").AsInt()
	}

	return 0
}

// AsString returns the current object as a string value
func (document Document) AsString() string {

	switch typed := document.value.(type) {

	case string:
		return typed

	case int:
		return strconv.Itoa(typed)

	case int64:
		return strconv.FormatInt(typed, 10)

	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)

	case bool:
		return strconv.FormatBool(typed)

	case map[string]any:
		return document.Get("id").AsString()

	case []any:
		return document.Get("id").AsString()
	}

	return ""
}

// AsTime returns the current object as a time value
func (document Document) AsTime() time.Time {

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
	}

	return time.Time{}
}

/******************************************
 * List-based Methods
 ******************************************/

// Head returns the first object in a slice.
// For all other document types, it returns the current document.
func (document Document) Head() Document {

	if slice, ok := document.value.([]any); ok {

		if len(slice) > 0 {

			return Document{
				value: slice[0],
				cache: document.cache,
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
				value: slice[1:],
				cache: document.cache,
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

func (document Document) sub(value any) Document {
	return Document{
		value: value,
		cache: document.cache,
	}
}
