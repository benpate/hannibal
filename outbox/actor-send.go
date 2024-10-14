package outbox

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

/******************************************
 * Sending Messages
 ******************************************/

// Send pushes a message onto the outbound queue.
// This currently uses the To and CC fields, but not BTo and BCC.
// https://www.w3.org/TR/activitypub/#delivery
func (actor *Actor) Send(message mapof.Any) {

	// Create a streams.Document from the message
	client := actor.getClient()
	document := streams.NewDocument(message, streams.WithClient(client))

	// Collect the list of recipients and other values required to send the message
	recipients := actor.getRecipients(document)
	uniquer := NewUniquer[string]()

	// Send the message to each recipient
	for recipientID := range recipients {

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

		// Don't send to duplicate addresses
		if uniquer.IsDuplicate(recipientID) {
			continue
		}

		// Use the recipientID to look up their inbox URL
		recipient, err := streams.NewDocument(recipientID, streams.WithClient(client)).Load()

		if err != nil {
			derp.Report(derp.Wrap(err, "hannibal.outbox.actor.Send", "Error loading recipient", recipientID))
			continue
		}

		inboxURL := recipient.Inbox().ID()

		if inboxURL == "" {
			log.Error().Msg("Recipient does not have an inbox")
			continue
		}

		// Send the request to the target Actor's inbox
		transaction := remote.Post(inboxURL).
			Accept(vocab.ContentTypeActivityPub).
			ContentType(vocab.ContentTypeActivityPub).
			With(SignRequest(*actor)).
			JSON(message).
			Queue(actor.queue)

		if canDebug() {
			transaction.With(options.Debug())
		}

		if err := transaction.Send(); err != nil {
			derp.Report(derp.Wrap(err, "hannibal.outbox.actor.Send", "Error sending ActivityPub request", inboxURL))
		}
	}
}

// getRecipients calculates the list of recipients for a given message
// and updates the message accordingly.
func (actor *Actor) getRecipients(message streams.Document) <-chan string {

	result := make(chan string)

	go func() {

		defer close(result)

		// Copy TO: field into recipients
		for to := range message.To().Channel() {
			log.Debug().Msg("getRecipients/To sending to: " + to.ID())
			result <- to.ID()
		}

		// Copy CC: field into recipients
		for cc := range message.CC().Channel() {
			log.Debug().Msg("getRecipients/CC sending to: " + cc.ID())
			result <- cc.ID()
		}

		// Copy Tag: field into recipients (Mentions only)
		for tag := range message.Object().Tag().Channel() {

			if tag.Type() != vocab.LinkTypeMention {
				continue
			}

			if href := tag.Href(); href != "" {
				log.Debug().Msg("getRecipients/Tag sending to: " + href)
				result <- href
			}
		}

		// Special rules for certain kinds of messages:
		switch message.Type() {

		// Accept activities are sent to the Actor of the original object
		// Return so that no other recipients are added.
		case vocab.ActivityTypeAccept:
			log.Debug().Msg("getRecipients/Accept sending to: " + message.Object().Actor().ID())
			result <- message.Object().Actor().ID()
			return

		// Follow messages are sent to the person being followed.
		// Return so that no other recipients are added.
		case vocab.ActivityTypeFollow:
			log.Debug().Msg("getRecipients/Follow sending to: " + message.Object().ID())
			result <- message.Object().ID()
			return

		// Delete and Undo messages are sent to all recipients of the original message
		case vocab.ActivityTypeDelete, vocab.ActivityTypeUndo:
			if object := message.Object(); object.NotNil() {
				for recipient := range actor.getRecipients(object) {
					log.Debug().Msg("getRecipients/Delete sending to: " + recipient)
					result <- recipient
				}
			}
			return

		// Like and Dislike messages are sent to the author of the original message
		case vocab.ActivityTypeAnnounce,
			vocab.ActivityTypeLike,
			vocab.ActivityTypeDislike:

			recipient := message.Object().Actor().ID()

			log.Trace().Msg("getRecipients/Announce/Like/Dislike sending to: " + recipient)
			result <- recipient

			// Don't return because we also want to tell
			// the world that we announce/like/dislike this thing
		}

		// Write Actors from inReplyTo properties
		if inReplyTo := message.InReplyTo(); inReplyTo.NotNil() {
			calcRecipients_inReplyTo(inReplyTo, result, 0)
		}

		// Finally, send the message to all of the Actor's Followers
		if actor.followers != nil {
			log.Debug().Msg("getRecipients/Follower: Scanning Followers...")
			for follower := range actor.followers {
				log.Debug().Msg("getRecipients/Follower sending to: " + message.Object().Actor().ID())
				result <- follower
			}
		} else {
			log.Debug().Msg("getRecipients/Follower: Followers channel is nil")
		}

		log.Debug().Msg("getRecipients/Done")
	}()

	// Return the channel
	return result
}

// calcRecipients_inReplyTo is a recursive function that searches for recipients
// in the "inReplyTo" property of a document, and all of its child `Object` documents.
func calcRecipients_inReplyTo(document streams.Document, result chan<- string, depth int) {

	// End recursion
	if document.IsNil() {
		return
	}

	// Maximum recursion depth.
	// TODO: Perhaps this should be a configurable value?
	if depth > 16 {
		return
	}

	// Add the actor of this document to the list of recipients
	if actor := document.Actor(); actor.NotNil() {
		log.Trace().Msg("calcRecipients_inReplyTo/Actor sending: " + actor.ID())
		result <- actor.ID()
	}

	// If this activity is "AtrributedTo" an actor, then add that actor to the list of recipients
	for attributedTo := document.AttributedTo(); attributedTo.NotNil(); attributedTo = attributedTo.Tail() {
		log.Trace().Msg("calcRecipients_inReplyTo/attributedTo sending: " + attributedTo.ID())
		result <- attributedTo.ID()
	}

	// Recursive search for "InReplyTo" fields. If this activity is a reply, then add the original author to the list of recipients
	for inReplyTo := document.InReplyTo(); inReplyTo.NotNil(); inReplyTo = inReplyTo.Tail() {
		log.Trace().Msg("calcRecipients_inReplyTo Recursing InReplyTo")
		calcRecipients_inReplyTo(inReplyTo, result, depth+1)
	}

	// Recursive search for replies in Object tree
	log.Trace().Msg("calcRecipients_inReplyTo Recursing Object")
	calcRecipients_inReplyTo(document.Object(), result, depth+1)
}
