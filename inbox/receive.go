package inbox

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
func ReceiveRequest(request *http.Request, client streams.Client, options ...Option) (document streams.Document, err error) {

	const location = "hannibal.pub.ReceiveRequest"

	config := NewReceiveConfig(options...)

	// Try to read the body from the request
	body, err := re.ReadRequestBody(request)

	if err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error reading body from request")
	}

	// Try to retrieve the object from the buffer
	document = streams.NilDocument(streams.WithClient(client))

	if err := json.Unmarshal(body, &document); err != nil {
		log.Err(err).Msg("Hannibal Inbox: Error Unmarshalling JSON")
		return streams.NilDocument(), derp.Wrap(err, location, "Error unmarshalling JSON body into ActivityPub document")
	}

	// Validate the document using injected Validators
	isValid := validateRequest(request, &document, config.Validators)

	if canDebug() && document.Type() != vocab.ActivityTypeDelete {
		requestBytes, _ := httputil.DumpRequest(request, true)

		fmt.Println("")
		fmt.Println("------------------------------------------")
		fmt.Println("HANNIBAL: Received Request:")
		fmt.Println(string(requestBytes))

		/*
			fmt.Println(request.Method + " " + request.URL.String() + " " + request.Proto)
			fmt.Println("Host: " + request.Host)
			for key, value := range request.Header {
				fmt.Println(key + ": " + strings.Join(value, ", "))
			}
			fmt.Println("")
			fmt.Println(string(body))
		*/
		fmt.Println("------------------------------------------")
		fmt.Println("")
	}

	// Log the request
	if !isValid {
		log.Trace().Err(err).Msg("Hannibal Inbox: Received document is not valid")
		return streams.NilDocument(), derp.NewUnauthorizedError(location, "Cannot validate received document", document.Value())
	}

	// Return the parsed document to the caller (vöïlä!)
	return document, nil
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
