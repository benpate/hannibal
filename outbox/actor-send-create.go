package outbox

import (
	"github.com/benpate/hannibal/datetime"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendCreate announces a newly created document (such as a "Note" or
// "Article") to the Actor's followers and the document's addressees.
func (actor *Actor) SendCreate(document streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendCreate: " + document.ID())
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeCreate,
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
