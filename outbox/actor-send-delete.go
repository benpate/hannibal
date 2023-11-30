package outbox

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
)

// SendDelete sends an "Delete" message to the recipient
// actor: The Actor that is sending the request
// activity: The activity that has been deleted
// recipient: The ActivityStreams profile of the message recipient
func (actor *Actor) SendDelete(activity streams.Document) {
	actor.SendActivity(vocab.ActivityTypeDelete, activity)
}
