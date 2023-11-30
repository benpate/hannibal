package outbox

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendFollow sends a "Follow" request to the recipient
// actor: The Actor that is sending the request
// followID: The unique ID of this request
// recipient: The ActivityStreams profile of the Actor that is being followed
func (actor *Actor) SendFollow(followID string, remoteActorID string) {

	// Build the ActivityStream "Follow" request
	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyID:     followID,
		vocab.PropertyType:   vocab.ActivityTypeFollow,
		vocab.PropertyActor:  actor.actorID,
		vocab.PropertyObject: remoteActorID,
	}

	// Send the request
	actor.Send(message)
}
