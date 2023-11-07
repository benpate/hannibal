package streams

import (
	"bytes"
	"encoding/gob"
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
	value      any
	httpHeader http.Header
	client     Client
}

// NewDocument creates a new Document object from a JSON-LD map[string]any
func NewDocument(value any, options ...DocumentOption) Document {

	result := Document{
		value:      normalize(value),
		httpHeader: make(http.Header),
		client:     NewDefaultClient(),
	}

	result.WithOptions(options...)
	return result
}

// NilDocument returns a new, empty Document.
func NilDocument(options ...DocumentOption) Document {
	return NewDocument(nil, options...)
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
	switch typed := document.value.(type) {
	case string:
		return typed == ""
	case map[string]any:
		return len(typed) == 0
	case []any:
		return len(typed) == 0
	case nil:
		return true
	default:
		return document.value == nil
	}
}

func (document Document) NotNil() bool {
	return !document.IsNil()
}

func (document Document) IsArray() bool {
	_, ok := document.value.([]any)
	return ok
}

func (document Document) IsMap() bool {
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

func (document Document) Clone() Document {

	result := Document{
		client:     document.client,
		httpHeader: document.httpHeader.Clone(),
	}

	switch typed := document.value.(type) {

	case string:
		result.value = typed
		return result

	case int:
		result.value = typed
		return result

	case int64:
		result.value = typed
		return result

	case float64:
		result.value = typed
		return result

	case map[string]any:
		result.value = map[string]any{}

	case mapof.Any:
		result.value = mapof.Any{}

	case []any:
		result.value = []any{}

	case sliceof.Any:
		result.value = sliceof.Any{}

	}

	buffer := new(bytes.Buffer)
	gob.NewEncoder(buffer).Encode(document.value)
	gob.NewDecoder(buffer).Decode(&result.value)

	return result
}

// Get returns a sub-property of the current document
func (document Document) Get(key string) Document {

	if result := document.get(key); !result.IsNil() {
		return result
	}

	return NilDocument()
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
			return document.sub(typed[0]).Get(key)
		}

	case sliceof.Any:
		if len(typed) > 0 {
			return document.sub(typed[0]).Get(key)
		}
	}

	return NilDocument()
}

// TODO: LOW: Add GetContext() method

/******************************************
 * Conversion Methods
 ******************************************/

// Array returns the array value of the current object
func (document Document) Slice() []any {
	return convert.SliceOfAny(document.value)
}

// SliceOfDocuments transforms the current object into a slice of separate
// Document objects, one for each value in the current document array.
func (document Document) SliceOfDocuments() sliceof.Object[Document] {
	values := document.Slice()
	result := make([]Document, 0, len(values))
	for _, value := range values {
		result = append(result, document.sub(value))
	}

	return result
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
func (document Document) Load(options ...any) (Document, error) {

	const location = "hannibal.streams.Document.Map"

	if document.IsNil() {
		return NilDocument(), nil
	}

	switch typed := document.value.(type) {

	case map[string]any:
		return document, nil

	case []any:
		return document.Head().Load(options...)

	case string:
		return document.getClient().Load(typed, options...)
	}

	return NilDocument(), derp.NewInternalError(location, "Document type is invalid", document.Value())
}

func (document Document) Map() map[string]any {

	switch typed := document.value.(type) {

	case map[string]any:
		return typed

	case []any:
		return document.Head().Map()

	case string:
		return map[string]any{
			vocab.PropertyID: typed,
		}

	default:
		return map[string]any{}
	}
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
 * Array-based Methods
 ******************************************/

// Len returns the length of the document.
// If the document is nil, then this method returns 0
// If the document is a slice, then this method returns the length of the slice
// Otherwise, this method returns 1
func (document Document) Len() int {

	if document.IsNil() {
		return 0
	}

	if slice, ok := document.value.([]any); ok {
		return len(slice)
	}

	if slice, ok := convert.SliceOfAnyOk(document.value); ok {
		return len(slice)
	}

	return 1
}

// At returns the document at the specified index.
// If this document is not a slice, then this method returns a nil document.
func (document Document) At(index int) Document {

	if slice, ok := document.value.([]any); ok {

		if index < len(slice) {
			return document.sub(slice[index])
		}
	}

	if slice, ok := convert.SliceOfAnyOk(document.value); ok {
		if index < len(slice) {
			return document.sub(slice[index])
		}
	}

	return NilDocument()
}

/******************************************
 * List-based Methods
 ******************************************/

// Head returns the first object in a slice.
// For all other document types, it returns the current document.
func (document Document) Head() Document {

	// Try it the easy way first
	if slice, ok := document.value.([]any); ok {
		if len(slice) > 0 {
			return document.sub(slice[0])
		}
	}

	// Try convert in case we have something ugly (like a primitive.A)
	if slice, ok := convert.SliceOfAnyOk(document.value); ok {
		if len(slice) > 0 {
			return document.sub(slice[0])
		}
	}

	return document
}

// Tail returns all records after the first in a slice.
// For all other document types, it returns a nil document.
func (document Document) Tail() Document {

	if slice, ok := document.value.([]any); ok {
		if len(slice) > 1 {
			return document.sub(slice[1:])
		}
	}

	// Try convert in case we have something ugly (like a primitive.A)
	if slice, ok := convert.SliceOfAnyOk(document.value); ok {
		if len(slice) > 1 {
			return document.sub(slice[1:])
		}
	}

	return NilDocument()
}

// IsEmpty return TRUE if the current object is empty
func (document Document) IsEmptyTail() bool {

	if slice, ok := document.value.([]any); ok {
		return len(slice) < 2
	}

	if slice, ok := convert.SliceOfAnyOk(document.value); ok {
		return len(slice) < 2
	}

	return true
}

/******************************************
 * TypeDetection
 ******************************************/

// IsTypeActor returns TRUE if this document represents any
// of the predefined actor types
func (document Document) IsTypeActor() bool {
	switch document.Type() {

	case
		vocab.ActorTypeApplication,
		vocab.ActorTypeGroup,
		vocab.ActorTypeOrganization,
		vocab.ActorTypePerson,
		vocab.ActorTypeService:

		return true
	}
	return false
}

// NotTypeActor returns TRUE if this document does NOT represent any
// of the predefined actor types
func (document Document) NotTypeActor() bool {
	return !document.IsTypeActor()
}

/******************************************
 * Helpers
 ******************************************/

// SetValue sets the value of this document to a new value.
func (document *Document) SetValue(value any) {
	document.value = value
}

// SetProperty sets an individual property within this document.
func (document *Document) SetProperty(property string, value any) {
	document.value = document.setProperty(document.value, property, value)
}

func (document *Document) setProperty(currentValue any, property string, value any) any {

	switch typed := currentValue.(type) {

	case map[string]any:
		typed[property] = value
		return typed

	case []any:
		if len(typed) == 0 {
			document.value = map[string]any{
				property: value,
			}
			return typed
		}

		firstItem := document.setProperty(typed[0], property, value)
		typed[0] = firstItem
		return typed

	case string:
		return map[string]any{
			vocab.PropertyID: typed,
			property:         value,
		}

	default:
		return map[string]any{
			vocab.PropertyID: typed,
			property:         value,
		}
	}
}

func (document *Document) WithOptions(options ...DocumentOption) {
	for _, option := range options {
		option(document)
	}
}

func (document *Document) getClient() Client {

	if document.client != nil {
		return document.client
	}

	return NewDefaultClient()
}

// sub returns a new Document with a new VALUE, all of the same OPTIONS as this original
func (document *Document) sub(value any) Document {
	return Document{
		value:      normalize(value),
		client:     document.client,
		httpHeader: document.httpHeader,
	}
}

func normalize(value any) any {

	switch typed := value.(type) {

	case mapof.Any:
		return map[string]any(typed)

	case sliceof.Any:
		return []any(typed)

	}

	return value
}
