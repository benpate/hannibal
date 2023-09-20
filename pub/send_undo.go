package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendUndoQueueTask creates a QueueTask that sends an "Undo" message to the recipient
func SendUndoQueueTask(actor Actor, activity mapof.Any, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendUndo(actor, activity, recipient)
	})
}

// SendUndo sends an "Undo" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been undone
// recipient: The ActivityStreams profile of the message recipient
func SendUndo(actor Actor, activity mapof.Any, recipient streams.Document) error {

	// Build the ActivityPub Message
	message := Undo(actor.ActorID, activity)

	// Send the message to the target
	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "hannibal.pub.PostUndoActivity", "Error sending Undo request", message)
	}

	return nil
}

func Undo(actorID string, activity mapof.Any) mapof.Any {

	delete(activity, "@context")

	// Build the ActivityPub Message
	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   vocab.ActivityTypeUndo,
		vocab.PropertyActor:  actorID,
		vocab.PropertyObject: activity,
	}

	return message
}
