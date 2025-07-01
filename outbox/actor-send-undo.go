package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendUndo sends an "Undo" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been undone
// recipient: The ActivityStream profile of the message recipient
func (actor *Actor) SendUndo(activity streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendUndo: " + activity.ID())
	}

	// Build the ActivityPub Message
	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeUndo,
		vocab.PropertyActor:     actor.ActorID,
		vocab.PropertyObject:    activity.Map(),
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	actor.Send(message, activity.RangeAddressees())
}
