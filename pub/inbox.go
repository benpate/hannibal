package activitypub

import (
	"bytes"
	"encoding/json"
	"net/http"

	vocab "github.com/EmissarySocial/emissary/protocols/activityvocabulary"
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
)

// ParseInboxRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ParseInboxRequest(request *http.Request) (string, mapof.Any, error) {

	const activityTypeError = "ERROR"

	const location = "activitypub.ParseInboxRequest"

	activity := mapof.NewAny()

	// RULE: Content-Type MUST be "application/activity+json"
	if request.Header.Get(ContentType) != ContentTypeActivityPub {
		return activityTypeError, activity, derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	}

	// TODO: Verify the request signature
	// RULE: Verify request signatures
	// verifier, err := httpsig.NewVerifier(request)

	// Try to read the body from the request
	bodyReader, err := request.GetBody()

	if err != nil {
		return activityTypeError, activity, derp.Wrap(err, location, "Error copying request body")
	}

	// Try to read the body into the buffer
	var bodyBuffer bytes.Buffer

	if _, err = bodyBuffer.ReadFrom(bodyReader); err != nil {
		return activityTypeError, activity, derp.Wrap(err, location, "Error reading body into buffer")
	}

	// Try to unmarshal the body from the buffer into a map.
	if err := json.Unmarshal(bodyBuffer.Bytes(), &activity); err != nil {
		return activityTypeError, activity, derp.Wrap(err, location, "Error unmarshalling body")
	}

	// First, assume that we have a fully defined activity
	if activityType := vocab.ValidateActivityType(activity.GetString("type")); activityType != vocab.Unknown {
		return activityType, activity, nil
	}

	// Otherwise, assume that we have an implicit "Create" activity
	if objectType := vocab.ValidateObjectType(activity.GetString("type")); objectType != vocab.Unknown {
		return vocab.ActivityTypeCreate, activity, derp.NewInternalError(location, "Implicit 'Create' activities are not yet implemented")
	}

	// Return the activity to the caller.
	return vocab.Unknown, activity, nil
}
