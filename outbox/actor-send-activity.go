package outbox

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// SendActivity wraps a document in a standard ActivityStream envelope and sends it to the target.
// actor: The Actor that is sending the request
// activityType: The type of activity to send (e.g. "Create", "Update", "Accept", etc)
// object: The object of the activity (e.g. the post that is being created, updated, etc)
// recipient: The ActivityStreams profile of the message recipient
func (actor *Actor) SendActivity(activityType string, object streams.Document) {

	log.Debug().Msg("outbox.Actor.SendActivity: " + activityType + ", objectId: " + object.ID())

	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   activityType,
		vocab.PropertyActor:  actor.actorID,
		vocab.PropertyObject: object.Map(streams.OptionStripContext),
	}

	actor.Send(message)
}
