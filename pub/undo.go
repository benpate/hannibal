package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

func SendUndoQueueTask(actor Actor, activity mapof.Any, targetURL string) QueueTask {
	return NewQueueTask(func() error {
		return SendUndo(actor, activity, targetURL)
	})
}

func SendUndo(actor Actor, activity mapof.Any, targetURL string) error {

	// Build the ActivityPub Message
	message := Undo(activity)
	message["actor"] = actor.ActorID

	// Send the message to the target
	if err := Send(actor, message, targetURL); err != nil {
		return derp.Wrap(err, "activitypub.PostUndoActivity", "Error sending Undo request")
	}

	return nil
}

// Undo creates an "Undo" activity for the specified activity,
// copying the "actor" value from the original activity.
func Undo(activity mapof.Any) mapof.Any {

	delete(activity, "@context")

	undo := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeUndo,
		"actor":    activity["actor"],
		"object":   activity,
	}

	return undo
}
