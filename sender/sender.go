// Package sender provides an object and queue consumer that
// can send outbound ActivityPub activities from an outbox
// to actor's inbox URLs. It automatically signs activities
// using the sending actor's private key.
package sender

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/ranges"
	"github.com/benpate/turbine/queue"
	"github.com/rs/zerolog/log"
)

// Sender manages delivery of outbound activities from the outbox,
// using the turbind Queue to deliver activities asynchronously.
// This object can be used on its own, or it can be embedded in the
// HTTPPoster object to automatically send activities when they are
// POST-ed to the outbox.
type Sender struct {
	queue   *queue.Queue
	locator Locator
}

// New returns a fully initialized Sender object
func New(locator Locator, q *queue.Queue) Sender {

	const location = "outbox2.Sender.New"

	// If we don't have an actual queue, then create one
	if q == nil {
		derp.Report(
			derp.Internal(
				location,
				"A background queue was not provided to the Outbox service.",
				"This configuration is okay for testing, but should not be used in production code.",
				"A proper background queue is required to recover from failed message deliveries.",
				"Yes, you. I'm talking to you. Don't do this in production.",
			),
		)

		// Make our own in-memory queue
		// This is bad bad bad for production systems
		q = queue.New()
	}

	// Return the Sender object
	return Sender{
		queue:   q,
		locator: locator,
	}
}

// Send queues a new task to deliver the provided activity to all recipients.
// IMPORTANT: The queue.Consumer in this package MUST be connected to a live
// queue process in order for outbound activities to be sent.
func (sender Sender) Send(activity mapof.Any) error {

	const location = "outbox2.Sender.Send"

	// NILCHECK: If the Outbox was not properly initialized with a queue, then report an error and return
	if sender.queue == nil {
		return derp.Internal(
			location,
			"Message cannot be sent because the background queue was not provided.",
			"This should never happen.",
		)
	}

	// Queue a new task to send this activity to all recipients.
	task := queue.NewTask(OutboxSendToAllRecipients, activity)

	if err := sender.queue.Publish(task); err != nil {
		return derp.Wrap(err, location, "Unable to enqueue outbound activity", activity)
	}

	// Success!
	return nil
}

// SendToAllRecipients sends a single ActivityPub activity from a the provided Actor to a single recipient's inbox URL.
func (sender *Sender) SendToAllRecipients(activity mapof.Any) queue.Result {

	const location = "hannibal.sender.SendToAllRecipients"

	// Locate the Actor that is sending this activity
	actorID := activity.GetString(vocab.PropertyActor)
	actor, err := sender.locator.Actor(actorID)

	if err != nil {
		return queue.Failure(derp.Wrap(err, location, "Unable to locate actor", activity))
	}

	// Use the Locator to resolve recipient URIs into inbox URLs
	recipients, err := getRecipients(sender.locator, activity)

	if err != nil {
		return queue.Error(derp.Wrap(err, location, "Unable to retrieve recipient addresses"))
	}

	// Remove duplicate recipient inbox URLs
	recipients = ranges.Unique(recipients)

	// Strip BCC and BTo fields before sending
	activity.Remove(vocab.PropertyBCC)
	activity.Remove(vocab.PropertyBTo)

	// Enqueue additional tasks to send this Activity to each recipient's inboxURL
	for recipient := range recipients {

		// Skip empty inbox URLs
		if recipient == "" {
			continue
		}

		log.Debug().Str("actorID", actorID).Str("recipient", recipient).Msg("Queueing outbound activity")

		task := queue.NewTask(OutboxSendToSingleRecipient, mapof.Any{
			"actor":    actor.ActorID(),
			"inbox":    recipient,
			"activity": activity,
		})

		if err := sender.queue.Publish(task); err != nil {
			return queue.Error(derp.Wrap(err, location, "Unable to enqueue outbound activity", "recipient", recipient))
		}
	}

	// Task Successed Successfully!
	return queue.Success()
}

// SendToSingleRecipient sends a single ActivityPub activity from a the provided Actor to a single recipient's inbox URL.
func (sender *Sender) SendToSingleRecipient(args mapof.Any) queue.Result {

	const location = "hannibal.sender.SendToSingleRecipient"

	// Collect arguments
	actorID := convert.String(args["actor"])
	inboxURL := convert.String(args["inbox"])
	activity := convert.MapOfAny(args["activity"])

	log.Debug().Str("actorID", actorID).Str("inboxURL", inboxURL).Msg("Sending outbound activity")

	// Locate the Actor that is sending this activity
	actor, err := sender.locator.Actor(actorID)

	if err != nil {
		return queue.Failure(derp.Wrap(err, location, "Unable to retrieve actor for outbound activity", "actorID: "+actorID))
	}

	// Prepare a transaction to send to target Actor's inbox
	transaction := remote.Post(inboxURL).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		With(signRequest(actor.PrivateKey())).
		JSON(activity)

	// Enable debugging (if requested)
	if canDebug() {
		transaction.With(options.Debug())
	}

	// Send the transaction to the recipient's inbox.
	// Errors will be handled by the asQueueResult() function in the queue.Consumer.
	if err := transaction.Send(); err != nil {

		// Special handling for HTTP 429 (Too Many Requests) error
		if tooManyRequests, retryDuration := derp.IsTooManyRequests(err); tooManyRequests {
			return queue.Requeue(retryDuration)
		}

		// If this is our fault then it can't be retried. Fail accordingly.
		if derp.IsClientError(err) {
			return queue.Failure(derp.Wrap(err, location, "Unable to send HTTP request (Client Error cannot be retried)"))
		}

		// Otherwise, it is a server error that can be retried by the standard queue mechanism.
		return queue.Error(derp.Wrap(err, location, "Unable to send HTTP request (Server Error can be retried)"))

	}

	// No error means the transaction was successful.  Woot woot!
	return queue.Success()
}
