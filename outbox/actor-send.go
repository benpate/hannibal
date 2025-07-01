package outbox

import (
	"iter"
	"net/url"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/mapof"
)

/******************************************
 * Sending Messages
 ******************************************/

// Send pushes a message onto the outbound queue, sending it to
// all recipients in the iterator.
// https://www.w3.org/TR/activitypub/#delivery
func (actor *Actor) Send(message mapof.Any, recipients ...iter.Seq[string]) {

	const location = "hannibal.outbox.actor.Send"

	// Send the message to each recipient
	for _, iterator := range recipients {

		for recipientID := range iterator {

			// Don't send to empty recipients
			if recipientID == "" {
				continue
			}

			// Don't send to the magic public recipient
			if recipientID == vocab.NamespaceActivityStreamsPublic {
				continue
			}

			// Don't send messages to myself
			if recipientID == actor.actorID {
				continue
			}

			if err := actor.SendOne(recipientID, message); err != nil {
				derp.Report(derp.Wrap(err, location, "Error sending message", recipientID))
			}
		}
	}
}

// SendOne sends a single message to a single recipient
func (actor *Actor) SendOne(recipientID string, message mapof.Any) error {

	const location = "hannibal.outbox.actor.SendOne"

	// Use the recipientID to look up their inbox URL
	recipient := streams.NewDocument(recipientID, streams.WithClient(actor.getClient()))
	recipient, err := recipient.Load()

	if err != nil {
		return derp.Wrap(err, location, "Error loading recipient", recipientID)
	}

	inboxURL := recipient.Inbox().ID()

	// RULE: InboxURL must be a valid URL
	inbox, err := url.Parse(inboxURL)

	if err != nil {
		return derp.Wrap(err, location, "Invalid Inbox URL", inboxURL)
	}

	// Prepare a transaction to send to target Actor's inbox
	transaction := remote.Post(inbox.String()).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		With(SignRequest(*actor)).
		JSON(message)

	if canDebug() {
		transaction.With(options.Debug())
	}

	// Send the transaction
	if err := transaction.Send(); err != nil {
		return derp.Wrap(err, location, "Error sending ActivityPub request", inboxURL)
	}

	return nil
}
