package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendDeleteQueueTask creates a QueueTask that sends an "Delete" message to the recipient
func SendDeleteQueueTask(actor Actor, activity streams.Document, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendDelete(actor, activity, recipient)
	})
}

// SendDelete sends an "Delete" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been deleted
// recipient: The ActivityStreams profile of the message recipient
func SendDelete(actor Actor, activity streams.Document, recipient streams.Document) error {

	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeDelete,
		"actor":    actor.ActorID,
		"object":   activity.Value(),
	}

	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
