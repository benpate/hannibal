package outbox

import (
	"time"

	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendAccept sends an "Announce" message to the recipient
// activity: The activity that is being announced
func (actor *Actor) SendAnnounce(announceID string, activity streams.Document) {

	log.Debug().Msg("outbox.Actor.SendAnnounce: " + announceID)

	message := mapof.Any{
		vocab.AtContext:         vocab.ContextTypeActivityStreams,
		vocab.PropertyType:      vocab.ActivityTypeAnnounce,
		vocab.PropertyID:        announceID,
		vocab.PropertyActor:     actor.actorID,
		vocab.PropertyObject:    activity.Map(streams.OptionStripContext),
		vocab.PropertyPublished: hannibal.TimeFormat(time.Now()),
	}

	actor.Send(message)
}
