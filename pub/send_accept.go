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

	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeAccept,
		"actor":    actor.ActorID,
		"object":   activity.Value(),
	}

	if err := Send(actor, message, activity.Actor()); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
