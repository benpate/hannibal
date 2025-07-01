package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendDislike sends an "Dislike" message to the recipient
// activity: The activity that is being announced
func (actor *Actor) SendDislike(dislikeID string, object streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendDislike: " + dislikeID)
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeDislike,
		vocab.PropertyID:        dislikeID,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    object.Map(),
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	actor.Send(message, actor.followers, object.RangeAddressees())
}
