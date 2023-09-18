package pub

import (
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

	document := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     activityType,
		"actor":    actor.ActorID,
		"object":   object,
	}

	return Send(actor, document, recipient)
}
