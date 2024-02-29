package inbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/re"
	"github.com/rs/zerolog/log"
)

// ReceiveRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ReceiveRequest(request *http.Request, client streams.Client) (document streams.Document, err error) {

	const location = "hannibal.pub.ReceiveRequest"

	// Try to read the body from the request
	body, err := re.ReadRequestBody(request)

	if err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error reading body from request")
	}

	// Debug if necessary
	if canDebug() {
		if canTrace() {
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
			log.Debug().Str("url", request.URL.String()).Msg("Hannibal Inbox: Received Request")
		}
	}

	// Try to retrieve the object from the buffer
	document = streams.NilDocument(streams.WithClient(client))

	if err := json.Unmarshal(body, &document); err != nil {
		log.Err(err).Msg("Hannibal Inbox: Error Unmarshalling JSON")
		return streams.NilDocument(), derp.Wrap(err, location, "Error unmarshalling JSON body into ActivityPub document")
	}

	// Validate the Actor and Public Key
	if err := validateRequest(request, document); err != nil {
		log.Err(err).Msg("Hannibal Inbox: Invalid Actor or Public Key")
		return streams.NilDocument(), derp.Wrap(err, location, "Invalid Actor or Public Key", document.Value())
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
