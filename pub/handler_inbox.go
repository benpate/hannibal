package pub

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/jsonld"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/davecgh/go-spew/spew"
)

// ParseInboxRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ParseInboxRequest(request *http.Request, client jsonld.Client) (activityType string, objectType string, reader jsonld.Reader, err error) {

	const activityTypeError = "ERROR"

	const location = "activitypub.ParseInboxRequest"

	activity := mapof.NewAny()

	// RULE: Content-Type MUST be "application/activity+json" or "application/ld+json"
	if !isActivityPubContentType(request.Header.Get(ContentType)) {
		return activityTypeError, vocab.Unknown, jsonld.NewZero(), derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	}

	// TODO: Verify the request signature
	// RULE: Verify request signatures
	// verifier, err := httpsig.NewVerifier(request)

	// Try to read the body from the request
	var bodyBuffer bytes.Buffer
	if _, err := bodyBuffer.ReadFrom(request.Body); err != nil {
		return activityTypeError, vocab.Unknown, jsonld.NewZero(), derp.Wrap(err, location, "Error reading body into buffer")
	}

	// Try to unmarshal the body from the buffer into a map.
	if err := json.Unmarshal(bodyBuffer.Bytes(), &activity); err != nil {
		return activityTypeError, vocab.Unknown, jsonld.NewZero(), derp.Wrap(err, location, "Error unmarshalling body")
	}

	spew.Dump("HandleInbox: received activity", activity)

	// First, assume that we have a fully defined activity
	if activityType := vocab.ValidateActivityType(activity.GetString("type")); activityType != vocab.Unknown {
		reader := client.NewReader(activity)
		objectType := reader.Get("object").Get("type").AsString()
		return activityType, objectType, reader, nil
	}

	// Otherwise, assume that we have an implicit "Create" activity
	if objectType := vocab.ValidateObjectType(activity.GetString("type")); objectType != vocab.Unknown {
		return vocab.ActivityTypeCreate, objectType, client.NewReader(activity), derp.NewInternalError(location, "Implicit 'Create' activities are not yet implemented")
	}

	// Return the activity to the caller.
	return vocab.Unknown, vocab.Unknown, client.NewReader(activity), nil
}
