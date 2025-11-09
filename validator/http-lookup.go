package validator

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/property"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
)

// HTTPLookup is a Validator that tries to retrieve the original document from the source server
type HTTPLookup struct{}

func NewHTTPLookup() HTTPLookup {
	return HTTPLookup{}
}

func (v HTTPLookup) Validate(request *http.Request, document *streams.Document) Result {

	// 	return ResultUnknown

	const location = "hannibal.validator.HTTPLookup"

	switch document.Type() {
	case vocab.ActivityTypeCreate, vocab.ActivityTypeUpdate:
	default:
		return ResultUnknown
	}

	// Get the ObjectID of the document
	objectID := document.Object().ID()

	if objectID == "" {
		return ResultInvalid
	}

	// Load the original document
	original, err := document.Client().Load(objectID)

	if err != nil {
		derp.Report(derp.Wrap(err, location, "Unable to load original document", objectID))
	}

	// Extract the value from the "original" retrieved document and replace it int the
	// document that was passed in
	value := original.Value()
	propertyValue := property.NewValue(value)
	document.SetValue(propertyValue)

	// Return in triumph
	return ResultValid
}
