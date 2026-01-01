package sender

import (
	"iter"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/turbine/queue"
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
func (sender *Sender) Send(activity mapof.Any) error {

	const location = "outbox2.Sender.Send"

	// NILCHECK: If the Outbox was not properly initialized with a queue, then report an error and return
	if sender.queue == nil {
		return derp.Internal(
			location,
			"Message cannot be sent because the background queue was not provided.",
			"This should never happen.",
		)
	}

	// Locate the Actor that is sending this activity
	actorID := activity.GetString(vocab.PropertyActor)
	actor, err := sender.locator.Actor(actorID)

	if err != nil {
		return derp.Wrap(err, location, "Unable to retrieve actor for outbound activity", "actorID", actorID)
	}

	// Use the Locator to resolve recipient URIs into inbox URLs
	recipients, err := sender.getRecipients(activity)

	if err != nil {
		return derp.Wrap(err, location, "Unable to retrieve of activity recipient addresses")
	}

	// Remove duplicate recipient inbox URLs
	uniquer := streams.NewUniquer[string]()
	recipients = uniquer.Range(recipients)

	// Strip BCC and BTo fields before sending
	activity.Remove(vocab.PropertyBCC)
	activity.Remove(vocab.PropertyBTo)

	// Queue up new tasks to send this Activity to each recipient's inboxURL
	for inboxURL := range recipients {

		queue.NewTask(OutboxSendActivity, mapof.Any{
			"actorID":  actor.ActorID(),
			"inboxURL": inboxURL,
			"activity": activity,
		})
	}

	// Success!
	return nil
}

func (sender *Sender) getRecipients(activity mapof.Any) (iter.Seq[string], error) {
	return getRecipients(sender.locator, activity)
}
