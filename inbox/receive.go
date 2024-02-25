package inbox

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/re"
	"github.com/rs/zerolog"
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

	// Logging
	log.Trace().Msg("------------------------------------")
	log.Debug().Str("inbox", request.Host+request.URL.String()).Msg("Hannibal Inbox: Activity Received")
	if canLog(zerolog.TraceLevel) {
		for key, value := range request.Header {
			log.Trace().Str(key, strings.Join(value, ", ")).Msg("Header")
		}
		log.Trace().Bytes("body", body).Msg("Body")
	}

	/*/ Debug if necessary
	if router.debug >= DebugLevelVerbose {
		fmt.Println("------------------------------------------")
		fmt.Println("HANNIBAL: Receiving Activity: " + request.URL.String())
		fmt.Println("Headers:")
		for key, value := range request.Header {
			fmt.Println(key + ": " + strings.Join(value, ", "))
		}
		fmt.Println("")
		fmt.Println("Body:")
		fmt.Println(string(body))
		fmt.Println("")
	}
	*/

	/* RULE: Content-Type MUST be "application/activity+json" or "application/ld+json"
	if !IsActivityPubContentType(request.Header.Get(vocab.ContentType)) {
		return streams.NilDocument(), derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	} */

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
	log.Debug().Str("id", document.ID()).Msg("Hannibal Inbox: Activity Parsed")
	if canLog(zerolog.TraceLevel) {
		rawJSON, _ := json.MarshalIndent(document.Value(), "", "  ")
		log.Trace().RawJSON("document", rawJSON).Send()
	}

	// Return the parsed document to the caller (vöïlä!)
	return document, nil
}
