package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendAnnounce sends an "Announce" activity for the provided object.
func (actor *Actor) SendAnnounce(announceID string, object streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendAnnounce: " + announceID)
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeAnnounce,
		vocab.PropertyID:        announceID,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    object.Map(),
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	actor.Send(message, actor.followers, object.RangeAddressees())
}
