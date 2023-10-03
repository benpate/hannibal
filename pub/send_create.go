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
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   vocab.ActivityTypeCreate,
		vocab.PropertyActor:  actor.ActorID,
		vocab.PropertyObject: activity,
	}

	if activityID, ok := activity[vocab.PropertyID].(string); ok {
		message[vocab.PropertyID] = activityID
	}

	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "hannibal.pub.PostAcceptActivity", "Error sending Accept request", message)
	}

	return nil
}
