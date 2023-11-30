package outbox

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendUndo sends an "Undo" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been undone
// recipient: The ActivityStreams profile of the message recipient
func (actor *Actor) SendUndo(activity streams.Document) {
	actor.Send(MakeUndo(actor.actorID, activity.Map(streams.OptionStripContext)))
}

func MakeUndo(actorID string, activity mapof.Any) mapof.Any {

	delete(activity, vocab.PropertyContext)

	// Build the ActivityPub Message
	return mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyType:   vocab.ActivityTypeUndo,
		vocab.PropertyActor:  actorID,
		vocab.PropertyObject: activity,
	}
}
