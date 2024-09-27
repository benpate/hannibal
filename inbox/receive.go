package inbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	// Debug if necessary
	if canDebug() {
		if canTrace() {
			if document.Type() != vocab.ActivityTypeDelete {
				fmt.Println("")
				fmt.Println("------------------------------------------")
				fmt.Println("HANNIBAL: Received Request:")
				fmt.Println(request.Method + " " + request.URL.String() + " " + request.Proto)
				fmt.Println("Host: " + request.Host)
				for key, value := range request.Header {
					fmt.Println(key + ": " + strings.Join(value, ", "))
				}
				fmt.Println("")
				fmt.Println(string(body))
				fmt.Println("------------------------------------------")
				fmt.Println("")
			} else {
				log.Debug().Str("object", document.Object().String()).Msg("Hannibal Inbox: Received Delete Activity")
			}
		} else {
			log.Debug().Str("url", request.URL.String()).Msg("Hannibal Inbox: Received Request")
		}
	}

	// Validate the document using injected Validators
	if !validateRequest(request, &document, config.Validators) {
		log.Err(err).Msg("Hannibal Inbox: Cannot validate received document")
		return streams.NilDocument(), derp.NewUnauthorizedError(location, "Cannot validate received document", document.Value())
	}

	// Logging
	if canDebug() {
		log.Debug().Str("id", document.ID()).Msg("Hannibal Inbox: Activity Parsed")
		if canTrace() {
			rawJSON, _ := json.MarshalIndent(document.Value(), "", "  ")
			log.Trace().RawJSON("document", rawJSON).Send()
		}
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
