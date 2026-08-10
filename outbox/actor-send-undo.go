package outbox

import (
	"github.com/benpate/hannibal/datetime"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendUndo announces that a previously sent activity has been undone, to that
// activity's original addressees.
func (actor *Actor) SendUndo(activity streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendUndo: " + activity.ID())
	}

	// Build the ActivityPub Message
	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeUndo,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    activity.Map(),
		vocab.PropertyPublished: datetime.Now(),
	}

	actor.Send(message, activity.RangeAddressees())
}
