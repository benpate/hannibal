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
	log.Debug().Msg("outbox.Actor.SendUndo: " + activity.ID())

	actor.Send(MakeUndo(actor.actorID, activity.Map()))
}

func MakeUndo(actorID string, activity mapof.Any) mapof.Any {

	context := activity[vocab.AtContext]
	delete(activity, vocab.AtContext)

	// Build the ActivityPub Message
	return mapof.Any{
		vocab.AtContext:         context,
		vocab.PropertyType:      vocab.ActivityTypeUndo,
		vocab.PropertyActor:     actorID,
		vocab.PropertyObject:    activity,
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}
}
