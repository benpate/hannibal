package outbox

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// SendCreate sends an "Create" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been created (such as a "Note" or "Article")
// recipient: The ActivityStreams profile of the message recipient
func (actor *Actor) SendCreate(activity mapof.Any) {

	message := mapof.Any{
		vocab.AtContext:      vocab.ContextTypeActivityStreams,
		vocab.PropertyID:     activity.GetString(vocab.PropertyID),
		vocab.PropertyType:   vocab.ActivityTypeCreate,
		vocab.PropertyActor:  actor.actorID,
		vocab.PropertyObject: activity,
	}

	actor.Send(message)
}
