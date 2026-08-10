package outbox

import (
	"github.com/benpate/hannibal/datetime"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendDelete announces that a document has been deleted, to the Actor's
// followers and the document's addressees.
func (actor *Actor) SendDelete(document streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendDelete: " + document.Object().ID())
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeDelete,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    document.Object().Map(),
		vocab.PropertyPublished: datetime.Now(),
	}

	actor.Send(message, document.RangeAddressees(), actor.followers)
}
