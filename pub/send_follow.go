package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendFollowQueueTask creates a QueueTask that sends a "Follow" request to the recipient
func SendFollowQueueTask(actor Actor, followID string, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendFollow(actor, followID, recipient)
	})
}

// SendFollow sends a "Follow" request to the recipient
// actor: The Actor that is sending the request
// followID: The unique ID of this request
// recipient: The ActivityStreams profile of the Actor that is being followed
func SendFollow(actor Actor, followID string, recipient streams.Document) error {

	// Build the ActivityStream "Follow" request
	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyID:     followID,
		vocab.PropertyType:   vocab.ActivityTypeFollow,
		vocab.PropertyActor:  actor.ActorID,
		vocab.PropertyObject: recipient.ID(),
	}

	// Send the request
	if err := Send(actor, message, recipient); err != nil {
		return derp.Wrap(err, "hannibal.pub.Follow", "Error sending Follow request", message)
	}

	return nil
}
