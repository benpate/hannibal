package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendFollowQueueTask creates a QueueTask that sends an "Accept" message to the recipient
func SendAcceptQueueTask(actor Actor, activity streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendAccept(actor, activity)
	})
}

// SendAccept sends an "Accept" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been accepted (likely a "Follow" request)
func SendAccept(actor Actor, activity streams.Document) error {

	object := activity.Map()
	delete(object, vocab.AtContext)

	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   vocab.ActivityTypeAccept,
		vocab.PropertyActor:  actor.ActorID,
		vocab.PropertyObject: object,
	}

	if err := Send(actor, message, activity.Actor()); err != nil {
		return derp.Wrap(err, "hannibal.pub.PostAcceptActivity", "Error sending Accept request", message)
	}

	return nil
}
