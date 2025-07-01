package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendDelete sends an "Delete" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been deleted
// recipient: The ActivityStream profile of the message recipient
func (actor *Actor) SendDelete(document streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendDelete: " + document.Object().ID())
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeDelete,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    document.Object().Map(),
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	actor.Send(message, document.RangeAddressees(), actor.followers)
}
