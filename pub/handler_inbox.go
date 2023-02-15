package pub

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/davecgh/go-spew/spew"
)

// ParseInboxRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ParseInboxRequest(request *http.Request, cache streams.Cache) (document streams.Document, err error) {

	const activityTypeError = "ERROR"

	const location = "activitypub.ParseInboxRequest"

	// RULE: Content-Type MUST be "application/activity+json" or "application/ld+json"
	if !isActivityPubContentType(request.Header.Get(vocab.ContentType)) {
		return streams.NilDocument(), derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	}

	// TODO: Verify the request signature
	// RULE: Verify request signatures
	// verifier, err := httpsig.NewVerifier(request)

	// Try to read the body from the request
	var bodyBuffer bytes.Buffer
	if _, err := bodyBuffer.ReadFrom(request.Body); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error reading body into buffer")
	}

	spew.Dump("ParseInboxRequest : RECEIVED ---------------------------", bodyBuffer.String(), "---------------------------")

	// Try to retrieve the object from the buffer
	document = streams.NewDocument(nil, cache)

	if err := json.Unmarshal(bodyBuffer.Bytes(), &document); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error unmarshalling JSON body into ActivityPub document")
	}

	messageType := document.Type()

	// First, assume that we have a fully defined activity
	if activityType := vocab.ValidateActivityType(messageType); activityType != vocab.Unknown {
		return document, nil
	}

	// Otherwise, assume that we have an implicit "Create" activity
	if objectType := vocab.ValidateObjectType(messageType); objectType != vocab.Unknown {
		// TODO: MEDIUM: Wrap original activity in a "Create" activity
		// TODO: MEDIUM: Can we get the Actor from the signed HTTP Request?
		return document, derp.NewInternalError(location, "Implicit 'Create' activities are not yet implemented")
	}

	// Return the activity to the caller.
	return streams.NilDocument(), derp.NewInternalError(location, "Unknown ActivityPub message type", document.Value())
}
