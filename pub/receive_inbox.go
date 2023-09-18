package pub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/re"
)

// ReceiveInboxRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ReceiveInboxRequest(request *http.Request, client streams.Client) (document streams.Document, err error) {

	const location = "activitypub.ReceiveInboxRequest"

	// Try to read the body from the request
	body, err := re.ReadRequestBody(request)

	if err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error reading body from request")
	}

	// Debug if necessary
	if packageDebugLevel >= DebugLevelVerbose {
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

	/* RULE: Content-Type MUST be "application/activity+json" or "application/ld+json"
	if !IsActivityPubContentType(request.Header.Get(vocab.ContentType)) {
		return streams.NilDocument(), derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	} */

	// Try to retrieve the object from the buffer
	document = streams.NilDocument(streams.WithClient(client))

	if err := json.Unmarshal(body, &document); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error unmarshalling JSON body into ActivityPub document")
	}

	// Validate the Actor and Public Key
	if err := validateRequest(request, document); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Request is invalid")
	}

	documentType := document.Type()

	// First, assume that we have a fully defined activity
	if activityType := vocab.ValidateActivityType(documentType); activityType != vocab.Unknown {
		return document, nil
	}

	// Otherwise, assume that we have an implicit "Create" activity
	if objectType := vocab.ValidateObjectType(documentType); objectType != vocab.Unknown {
		// TODO: MEDIUM: Wrap original activity in a "Create" activity
		// TODO: MEDIUM: Can we get the Actor from the signed HTTP Request?
		return document, derp.NewInternalError(location, "Implicit 'Create' activities are not yet implemented")
	}

	// Return the activity to the caller.
	return streams.NilDocument(), derp.NewInternalError(location, "Unknown ActivityPub message type", document.Value())
}
