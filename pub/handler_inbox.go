package pub

import (
	"bytes"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/jsonld"
	"github.com/benpate/hannibal/vocab"
)

// ParseInboxRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ParseInboxRequest(request *http.Request, client jsonld.Client) (reader jsonld.Reader, err error) {

	const activityTypeError = "ERROR"

	const location = "activitypub.ParseInboxRequest"

	// RULE: Content-Type MUST be "application/activity+json" or "application/ld+json"
	if !isActivityPubContentType(request.Header.Get(ContentType)) {
		return jsonld.NilReader(), derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	}

	// TODO: Verify the request signature
	// RULE: Verify request signatures
	// verifier, err := httpsig.NewVerifier(request)

	// Try to read the body from the request
	var bodyBuffer bytes.Buffer
	if _, err := bodyBuffer.ReadFrom(request.Body); err != nil {
		return jsonld.NilReader(), derp.Wrap(err, location, "Error reading body into buffer")
	}

	// Try to unmarshal the body from the buffer into a new JSON-LD reader
	reader = client.UnmarshalReader(bodyBuffer.Bytes())

	messageType := reader.Type()

	// First, assume that we have a fully defined activity
	if activityType := vocab.ValidateActivityType(messageType); activityType != vocab.Unknown {
		return reader, nil
	}

	// Otherwise, assume that we have an implicit "Create" activity
	if objectType := vocab.ValidateObjectType(messageType); objectType != vocab.Unknown {
		// TODO: MEDIUM: Wrap original activity in a "Create" activity
		// TODO: MEDIUM: Can we get the Actor from the signed HTTP Request?
		return reader, derp.NewInternalError(location, "Implicit 'Create' activities are not yet implemented")
	}

	// Return the activity to the caller.
	return reader, nil
}
