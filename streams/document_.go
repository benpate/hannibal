package streams

import (
	"html"
	"net/http"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/property"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/sliceof"
	"github.com/microcosm-cc/bluemonday"
)

// Document represents a single ActivityStream document
// or document fragment.  Due to the flexibility of ActivityStreams
// (and JSON-LD), this may be a data structure such as a
// `map[string]any`, `[]any`, or a primitive type, like a
// `string`, `float`, `int` or `bool`.
type Document struct {
	value      property.Value
	statistics Statistics
	httpHeader http.Header
	client     Client
}

// NewDocument creates a new Document object from a JSON-LD map[string]any
func NewDocument(value any, options ...DocumentOption) Document {

	result := Document{
		value:      property.NewValue(value),
		statistics: NewStatistics(),
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

func (document Document) IsBool() bool {
	return property.IsBool(document.value)
}

func (document Document) IsInt() bool {
	return property.IsInt(document.value)
}

func (document Document) IsInt64() bool {
	return property.IsInt64(document.value)
}

func (document Document) IsFloat() bool {
	return property.IsFloat(document.value)
}

func (document Document) IsNil() bool {
	return document.value.IsNil()
}

func (document Document) IsMap() bool {
	return property.IsMap(document.value)
}

func (document Document) IsSlice() bool {
	return property.IsSlice(document.value)
}

func (document Document) IsString() bool {
	return property.IsString(document.value)
}

func (document Document) NotNil() bool {
	return !document.IsNil()
}

/******************************************
 * Getter Methods
 ******************************************/

// Value returns the generic data stored in this Document
func (document Document) Value() any {
	return document.value.Raw()
}

func (document Document) Clone() Document {

	return Document{
		client:     document.client,
		httpHeader: document.httpHeader.Clone(),
		value:      document.value.Clone(),
	}
}

// Get returns a sub-property of the current document
func (document Document) Get(key string) Document {

	// Special handling for string values
	if document.IsString() {

		// Individual values are assumed to be a document ID.
		// So if te ID property was requested, then just return it
		if key == vocab.PropertyID {
			return document
		}

		// All other properties require a Load from the Interweb.
		// Update this document with the loaded result.
		if loaded, err := document.Load(); err == nil {
			document.value = loaded.value
		}
	}

	// Retrieve the value from the property.Value
	if value := document.value.Get(key); !value.IsNil() {
		return document.sub(value)
	}

	// Nil document if the property doesn't exist
	return NilDocument()
}

// TODO: LOW: Add GetContext() method

/******************************************
 * Conversion Methods
 ******************************************/

// Array returns the array value of the current object
func (document Document) Slice() []any {
	return convert.SliceOfAny(document.value.Raw())
}

// SliceOfDocuments transforms the current object into a slice of separate
// Document objects, one for each value in the current document array.
func (document Document) SliceOfDocuments() sliceof.Object[Document] {
	values := document.Slice()
	result := make([]Document, 0, len(values))
	for _, value := range values {
		result = append(result, document.sub(property.NewValue(value)))
	}

	return result
}

// Bool returns the current object as a floating-point value
func (document Document) Bool() bool {

	if getter, ok := document.value.(property.BoolGetter); ok {
		return getter.Bool()
	}

	return convert.Bool(document.value.Head().Raw())
}

// Float returns the current object as an integer value
func (document Document) Float() float64 {

	if getter, ok := document.value.(property.FloatGetter); ok {
		return getter.Float()
	}

	return convert.Float(document.value.Head().Raw())
}

// Int returns the current object as an integer value
func (document Document) Int() int {

	if getter, ok := document.value.(property.IntGetter); ok {
		return getter.Int()
	}

	return convert.Int(document.value.Head().Raw())
}

// Load retrieves a JSON-LD document from its remote server
func (document Document) Load(options ...any) (Document, error) {

	// Get the document ID
	documentID := document.ID()

	// Guarantee that we have an actual value
	if documentID == "" {
		return NilDocument(), nil
	}

	// Try to load the document from the Interwebs
	result, err := document.getClient().Load(documentID, options...)

	if err != nil {
		return result, derp.Wrap(err, "hannibal.streams.document.Load", "Error loading document by ID", document.Value())
	}

	// Success??
	return result, nil
}

// MustLoad retrieves a JSON-LD document from its remote server.
// It silently reports errors, but does not return them.
func (document Document) MustLoad(options ...any) Document {
	result, err := document.Load(options...)
	if err != nil {
		derp.Report(derp.Wrap(err, "hannibal.streams.document.MustLoad", "Error loading document"))
	}
	return result
}

// LoadLink loads a new JSON-LD document from a link or ID string.
// If the current document has already been loaded (because it's a map)
// then it is returned as-is.
func (document Document) LoadLink(options ...any) Document {

	// If this document is a string, then assume it's
	// an ID and load it from the Intertubes.
	if document.IsString() {
		return document.MustLoad(options...)
	}

	// Nothing to load. We already have a map.
	return document
}

func (document Document) Map(options ...string) map[string]any {

	// Create an empty result map
	result := make(map[string]any)

	// Traverse slices, if necessary
	value := document.value.Head()

	if getter, ok := value.(property.MapGetter); ok {
		result = getter.Map()
	} else if getter, ok := value.(property.IsStringer); ok {
		result[vocab.PropertyID] = getter.String()
	}

	// Apply optional filters
	for _, option := range options {
		switch option {

		case OptionStripContext:
			delete(result, vocab.AtContext)

		case OptionStripRecipients:
			delete(result, vocab.PropertyTo)
			delete(result, vocab.PropertyBTo)
			delete(result, vocab.PropertyCC)
			delete(result, vocab.PropertyBCC)
		}
	}

	return result

}

func (document Document) MapKeys() []string {

	if mapper, ok := document.value.(property.IsMapper); ok {
		return mapper.MapKeys()
	}

	return []string{}
}

// String returns the current object as a pure string (no HTML).
// This value is filtered by blueMonday, so it is safe to use in HTML.
func (document Document) String() string {
	result := document.rawString()
	result = bluemonday.StrictPolicy().Sanitize(result)
	result = html.UnescapeString(result)
	return result
}

// StringHTML returns the current object as an HTML string.
// This value is filtered by blueMonday, so it is safe to use in HTML.
func (document Document) HTMLString() string {
	result := document.rawString()
	return bluemonday.UGCPolicy().Sanitize(result)
}

// String returns the current object as a string value
func (document Document) rawString() string {

	if getter, ok := document.value.(property.IsStringer); ok {
		return getter.String()
	}

	return convert.String(document.value.Head().Raw())
}

// Time returns the current object as a time value
func (document Document) Time() time.Time {

	if getter, ok := document.value.(property.TimeGetter); ok {
		return getter.Time()
	}

	return convert.Time(document.value.Head().Raw())
}

/******************************************
 * Array-based Iterators
 ******************************************/

// Len returns the length of the document.
// If the document is nil, then this method returns 0
// If the document is a slice, then this method returns the length of the slice
// Otherwise, this method returns 1
func (document Document) Len() int {
	return document.value.Len()
}

/******************************************
 * List-based Iterators
 ******************************************/

// Head returns the first object in a slice.
// For all other document types, it returns the current document.
func (document Document) Head() Document {
	return document.sub(document.value.Head())
}

// Tail returns all records after the first in a slice.
// For all other document types, it returns a nil document.
func (document Document) Tail() Document {
	return document.sub(document.value.Tail())
}

// IsEmpty return TRUE if the current object is empty
func (document Document) IsEmptyTail() bool {
	return document.value.Len() < 2
}

/******************************************
 * Channel Iterators
 ******************************************/

// Channel returns a channel that iterates over all of the sub-documents
// in the current document.
func (document Document) Channel() <-chan Document {

	result := make(chan Document)

	go func() {
		defer close(result)

		for document.NotNil() {
			result <- document.Head()
			document = document.Tail()
		}
	}()

	return result
}

/******************************************
 * Helpers
 ******************************************/

// Client returns the HTTP client used for this document
func (document *Document) Client() Client {
	return document.client
}

// SetValue sets the value of this document to a new value.
func (document *Document) SetValue(value property.Value) {
	document.value = value
}

// SetProperty sets an individual property within this document.
func (document *Document) SetProperty(property string, value any) {
	document.value = document.value.Set(property, value)
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
func (document *Document) sub(value property.Value) Document {
	return Document{
		value:      value,
		client:     document.client,
		httpHeader: document.httpHeader,
	}
}
