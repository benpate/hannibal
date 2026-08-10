package outbox

import (
	"github.com/benpate/hannibal/datetime"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendFollow sends a "Follow" request to the designated remote Actor.
func (actor *Actor) SendFollow(followID string, remoteActorID string) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendFollow: " + followID)
	}

	// Build the ActivityStream "Follow" request
	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyID:        followID,
		vocab.PropertyType:      vocab.ActivityTypeFollow,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    remoteActorID,
		vocab.PropertyPublished: datetime.Now(),
	}

	// Send the request
	actor.Send(message, makeIterator(remoteActorID))
}
