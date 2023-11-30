package outbox

import (
	"encoding/json"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/******************************************
 * Sending Messages
 ******************************************/

// Send pushes a message onto the outbound queue.
// This currently uses the To and CC fields, but not BTo and BCC.
// https://www.w3.org/TR/activitypub/#delivery
func (actor *Actor) Send(message mapof.Any) {

	const location = "hannibal.outbox.Actor.Send"

	logger := log.With().Str("loc", location).Logger()

	if canLog(zerolog.DebugLevel) {
		logger.Debug().Msg("Sending message...")
		rawJSON, _ := json.MarshalIndent(message, "", "  ")
		logger.Debug().Msg(string(rawJSON))
	}

	// Create a streams.Document from the message
	document := streams.NewDocument(message, streams.WithClient(actor.getClient()))

	// Collect the list of recipients and other values required to send the message
	recipients := actor.getRecipients(document)
	client := actor.getClient()
	queue := actor.getQueue()
	uniquer := NewUniquer[string]()

	// Send the message to each recipient
	for recipient := range recipients {

		logger.Trace().Msg("Found Recipient: " + recipient)

		// Don't send to empty recipients
		if recipient == "" {
			logger.Trace().Msg("Empty recipient.")
			continue
		}

		// Don't send to the magic public recipient
		if recipient == vocab.NamespaceActivityStreamsPublic {
			logger.Trace().Msg("Public recipient. Do not deliver to the public namespace.")
			continue
		}

		// Don't send to duplicate addresses
		if uniquer.IsDuplicate(recipient) {
			logger.Trace().Msg("Duplicate recipient.")
			continue
		}

		// Make a copy of the message, individualized for this recipient,
		// and adding the recipient in the To field.
		messageMap := document.Map(streams.OptionStripRecipients)
		messageMap[vocab.PropertyTo] = recipient

		logger.Debug().Str("recipient", recipient).Msg("Queuing SendTask...")

		// Send the message to the recipient
		recipientDocument := streams.NewDocument(recipient, streams.WithClient(client))
		task := NewSendTask(*actor, messageMap, recipientDocument)
		queue.Push(task)
	}
}

// getRecipients calculates the list of recipients for a given message
// and updates the message accordingly.
func (actor *Actor) getRecipients(message streams.Document) <-chan string {

	result := make(chan string)

	go func() {

		defer close(result)

		// Copy TO: field into recipients
		for to := message.To(); to.NotNil(); to.Tail() {
			result <- to.Head().ID()
		}

		// Copy CC: field into recipients
		for cc := message.CC(); cc.NotNil(); cc.Tail() {
			result <- cc.Head().ID()
		}

		// Special rules for certain kinds of messages:
		switch message.Type() {

		// Accept activities are sent to the Actor of the original object
		// Return so that no other recipients are added.
		case vocab.ActivityTypeAccept:
			result <- message.Object().Actor().ID()
			return

		// Follow messages are sent to the person being followed.
		// Return so that no other recipients are added.
		case vocab.ActivityTypeFollow:
			result <- message.Object().ID()
			return

		// Like and Dislike messages are sent to the author of the original message
		case vocab.ActivityTypeAnnounce,
			vocab.ActivityTypeLike,
			vocab.ActivityTypeDislike:

			result <- message.Object().Actor().ID()

			// Don't return because we also want to tell
			// the world that we announce/like/dislike this thing
		}

		// Write Actors from inReplyTo properties
		calcRecipients_inReplyTo(message, result, 0)

		// Finally, send the message to all of the Actor's Followers
		if actor.followers != nil {
			for follower := range actor.followers {
				result <- follower
			}
		}
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
		result <- actor.ID()
	}

	// If this activity is "AtrributedTo" an actor, then add that actor to the list of recipients
	for attributedTo := document.AttributedTo(); attributedTo.NotNil(); attributedTo = attributedTo.Tail() {
		result <- attributedTo.ID()
	}

	// Recursive search for "InReplyTo" fields. If this activity is a reply, then add the original author to the list of recipients
	for inReplyTo := document.InReplyTo(); inReplyTo.NotNil(); inReplyTo = inReplyTo.Tail() {
		calcRecipients_inReplyTo(inReplyTo, result, depth+1)
	}

	// Recursive search for replies in Object tree
	calcRecipients_inReplyTo(document.Object(), result, depth+1)
}