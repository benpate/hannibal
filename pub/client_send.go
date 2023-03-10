package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

// Send sends an ActivityStream to a remote ActivityPub service
// actor: The Actor that is sending the request
// document: The ActivityStream that is being sent
// targetID: The ID of the Actor that will receive the request
func Send(actor Actor, document mapof.Any, targetID string) error {

	const location = "hannibal.pub.Send"

	// Try to get the source profile that we're going to follow
	target, err := GetProfile(targetID)

	if err != nil {
		return derp.Wrap(err, location, "Error getting source profile", targetID)
	}

	// Try to get the actor's inbox from the actor ActivityStream.
	// TODO: LOW: Is there a better / more reliable way to do this?
	inbox := target.GetString("inbox")

	if inbox == "" {
		return derp.NewInternalError(location, "Unable to find 'inbox' in target profile", targetID, target)
	}

	// Send the request to the target Actor's inbox
	transaction := remote.Post(inbox).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		Use(RequestSignature(actor)).
		JSON(document)

	if err := transaction.Send(); err != nil {
		return derp.Wrap(err, location, "Error sending Follow request", inbox)
	}

	// Done!
	return nil
}

// SendActivity wraps a document in a standard ActivityStream envelope and sends it to the target.
func SendActivity(actor Actor, activityType string, object any, targetID string) error {

	document := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     activityType,
		"actor":    actor.ActorID,
		"object":   object,
	}

	return Send(actor, document, targetID)
}
