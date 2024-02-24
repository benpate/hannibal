package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendFollow sends a "Follow" request to the recipient
// actor: The Actor that is sending the request
// followID: The unique ID of this request
// recipient: The ActivityStream profile of the Actor that is being followed
func (actor *Actor) SendFollow(followID string, remoteActorID string) {

	log.Debug().Msg("outbox.Actor.SendFollow: " + followID)

	// Build the ActivityStream "Follow" request
	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyID:        followID,
		vocab.PropertyType:      vocab.ActivityTypeFollow,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    remoteActorID,
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	// Send the request
	actor.Send(message)
}
