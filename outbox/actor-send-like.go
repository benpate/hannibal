package outbox

import (
	"github.com/benpate/hannibal/datetime"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendLike announces that the Actor has liked an object, to the Actor's
// followers and the object's addressees.
func (actor *Actor) SendLike(likeID string, object streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendLike: " + likeID)
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeLike,
		vocab.PropertyID:        likeID,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    object.Map(),
		vocab.PropertyPublished: datetime.Now(),
	}

	actor.Send(message, actor.followers, object.RangeAddressees())
}
