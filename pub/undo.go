package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

func SendUndoQueueTask(actor Actor, activity streams.Document, targetURL string) QueueTask {
	return NewQueueTask(func() error {
		return SendUndo(actor, activity, targetURL)
	})
}

func SendUndo(actor Actor, activity streams.Document, targetURL string) error {

	// RULE: Guarantee that the actor is the same one that created the original activity
	if actor.ActorID != activity.ActorID() {
		return derp.NewInternalError("activitypub.SendUndo", "Cannot undo an activity that was not created by this actor", nil)
	}

	// Build the ActivityPub Message
	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeUndo,
		"actor":    activity.ActorID(),
		"object":   activity.Value(),
	}

	// Send the message to the target
	if err := Send(actor, message, targetURL); err != nil {
		return derp.Wrap(err, "activitypub.PostUndoActivity", "Error sending Undo request")
	}

	return nil
}
