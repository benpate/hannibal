package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendUpdateQueueTask creates a QueueTask that sends an "Update" message to the recipient
func SendUpdateQueueTask(actor Actor, activity mapof.Any, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendUpdate(actor, activity, recipient)
	})
}

// SendUpdate sends an "Update" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been updated
// recipient: The ActivityStreams profile of the message recipient
func SendUpdate(actor Actor, activity mapof.Any, recipient streams.Document) error {

	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeUpdate,
		"actor":    actor.ActorID,
		"object":   activity,
	}

	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
