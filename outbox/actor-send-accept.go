package outbox

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendAccept sends an "Accept" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been accepted (likely a "Follow" request)
func (actor *Actor) SendAccept(acceptID string, activity streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendAccept: " + acceptID)
	}

	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyID:     acceptID,
		vocab.PropertyType:   vocab.ActivityTypeAccept,
		vocab.PropertyActor:  actor.actorID,
		vocab.PropertyObject: activity.Map(),
	}

	recipients := activity.Actor().RangeIDs()

	actor.Send(message, recipients)
}
