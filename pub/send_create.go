package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendCreateQueueTask creates a QueueTask that sends an "Create" message to the recipient
func SendCreateQueueTask(actor Actor, activity mapof.Any, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendCreate(actor, activity, recipient)
	})
}

// SendCreate sends an "Create" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been created (such as a "Note" or "Article")
// recipient: The ActivityStreams profile of the message recipient
func SendCreate(actor Actor, activity mapof.Any, recipient streams.Document) error {

	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeCreate,
		"actor":    actor.ActorID,
		"object":   activity,
	}

	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
