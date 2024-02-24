package outbox

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
)

// SendUpdate sends an "Update" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been updated
// recipient: The ActivityStream profile of the message recipient
func (actor *Actor) SendUpdate(activity streams.Document) {
	actor.SendActivity(vocab.ActivityTypeUpdate, activity)
}
