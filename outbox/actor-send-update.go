package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendUpdate sends an "Update" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been updated
// recipient: The ActivityStream profile of the message recipient
func (actor *Actor) SendUpdate(document streams.Document) {

	if canDebug() {
		log.Debug().Msg("outbox.Actor.SendUpdate: " + document.ID())
	}

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeUpdate,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    document.Map(),
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	actor.Send(
		message,
		document.RangeAddressees(),
		document.RangeInReplyTo(),
		actor.followers,
	)
}
