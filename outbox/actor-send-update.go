package outbox

import (
	"github.com/benpate/hannibal/datetime"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendUpdate announces that a document has been updated, to the Actor's
// followers and the document's addressees.
func (actor *Actor) SendUpdate(document streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendUpdate: " + document.ID())
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeUpdate,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    document.Map(),
		vocab.PropertyPublished: datetime.Now(),
	}

	actor.Send(
		message,
		document.RangeAddressees(),
		document.RangeInReplyTo(),
		actor.followers,
	)
}
