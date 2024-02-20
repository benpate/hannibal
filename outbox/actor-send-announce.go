package outbox

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendAccept sends an "Announce" message to the recipient
// activity: The activity that is being announced
func (actor *Actor) SendAnnounce(activity streams.Document) {

	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   vocab.ActivityTypeAnnounce,
		vocab.PropertyActor:  actor.actorID,
		vocab.PropertyObject: activity.Map(streams.OptionStripContext),
	}

	actor.Send(message)
}
