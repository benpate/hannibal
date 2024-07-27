package inbox

import (
	"net/http"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/validator"
)

// Validator interface wraps the Validate method, which identifies whether a document
// received in an actor's inbox is valid or not.  Multiple validators can be stacked
// to validate a document, so if one validator returns `false`, the document is not
// necessary invalid.  It just can't be validated by this one validator.
type Validator interface {

	// Validate checks incoming HTTP requests for validity.  If a document is
	// valid, it returns `ResultValid`.  If the Validator cannot validate this
	// document, it returns `ResultUnknown`. If the Validator can say with
	// certainty that the document is invalid, it returns `ResultInvalid`.
	Validate(*http.Request, *streams.Document) validator.Result
}
