package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendFollowQueueTask creates a QueueTask that sends a "Follow" request to the target Actor
func SendFollowQueueTask(actor Actor, followID string, targetID string) QueueTask {
	return NewQueueTask(func() error {
		return SendFollow(actor, followID, targetID)
	})
}

// SendFollow sends a "Follow" request to the target Actor
// actor: The Actor that is sending the request
// followID: The unique ID of this request
// targetID: The ID of the Actor that is being followed
func SendFollow(actor Actor, followID string, targetID string) error {

	// Build the ActivityStream "Follow" request
	activity := mapof.Any{
		"@context": vocab.ContentTypeActivityPub,
		"id":       followID,
		"type":     vocab.ActivityTypeFollow,
		"actor":    actor.ActorID,
		"object":   targetID,
	}

	// Send the request
	if err := Post(actor, activity, targetID); err != nil {
		return derp.Wrap(err, "activitypub.Follow", "Error sending Follow request")
	}

	return nil
}
