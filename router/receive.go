package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/validator"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/re"
	"github.com/rs/zerolog/log"
)

// ReceiveRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ReceiveRequest(request *http.Request, client streams.Client, options ...Option) (activity streams.Document, err error) {

	const location = "hannibal.router.ReceiveRequest"

	config := NewReceiveConfig(options...)

	// Try to read the body from the request
	body, err := re.ReadRequestBody(request)

	if err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Unable to read body from request")
	}

	// Try to retrieve the object from the buffer
	activity = streams.NilDocument(streams.WithClient(client))

	if err := json.Unmarshal(body, &activity); err != nil {
		log.Err(err).Msg("Hannibal Router: Error Unmarshalling JSON")
		return streams.NilDocument(), derp.Wrap(err, location, "Error unmarshalling JSON body into ActivityPub activity")
	}

	// Log the request
	if canDebug() && activity.Type() != vocab.ActivityTypeDelete {
		requestBytes, _ := httputil.DumpRequest(request, true)

		fmt.Println("")
		fmt.Println("Begin: Hannibal ReceiveRequest -----------")
		fmt.Println(string(requestBytes))
		fmt.Println("------------------------------------------")
		fmt.Println("")
	}

	// Validate the activity using injected Validators
	if isValid := validateRequest(request, &activity, config.Validators); !isValid {
		log.Trace().Msg("Hannibal Router: Received activity is not valid")
		return streams.NilDocument(), derp.Unauthorized(location, "Cannot validate received activity", activity.Value())
	}

	// Return the parsed activity to the caller (vöïlä!)
	return activity, nil
}

func validateRequest(request *http.Request, document *streams.Document, validators []Validator) bool {

	// Run each validator
	for _, v := range validators {

		switch v.Validate(request, document) {

		case validator.ResultInvalid:
			return false

		case validator.ResultValid:
			return true

		}

		// Fall through "ResultUnknown"
		// means continue the loop.
	}

	// If no validators can actually validate the document, then validation fails
	return false
}
