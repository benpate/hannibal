package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendActivityQueueTask creates a QueueTask that sends a message to the recipient
func SendActivityQueueTask(actor Actor, activityType string, object any, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendActivity(actor, activityType, object, recipient)
	})
}

// SendActivity wraps a document in a standard ActivityStream envelope and sends it to the target.
// actor: The Actor that is sending the request
// activityType: The type of activity to send (e.g. "Create", "Update", "Accept", etc)
// object: The object of the activity (e.g. the post that is being created, updated, etc)
// recipient: The ActivityStreams profile of the message recipient
func SendActivity(actor Actor, activityType string, object any, recipient streams.Document) error {

	if objectMap, ok := object.(map[string]any); ok {
		delete(objectMap, vocab.AtContext)
	}

	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   activityType,
		vocab.PropertyActor:  actor.ActorID,
		vocab.PropertyObject: object,
	}

	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "hannibal.pub.SendActivity", "Error sending Activity", message)
	}

	return nil
}
