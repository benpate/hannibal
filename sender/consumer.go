package sender

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/turbine/queue"
)

// Consumer returns a turbine queue.Consumer that processes
// outbound ActivityPub actitities for this outbox.
func Consumer(locator Locator) queue.Consumer {

	return func(name string, args map[string]any) queue.Result {

		switch name {

		// Catalog all recipients and queue individual send tasks
		case OutboxSendActivity:
			return sendActivity(locator, args)

		// Send an activity to a single recipient
		case OutboxSendActivity_SingleRecipient:
			return sendActivity_SingleRecipient(locator, args)
		}

		// All other task names are left for other consumers.
		return queue.Ignored()
	}
}

// sendActivity sends a single ActivityPub activity from a the provided Actor to a single recipient's inbox URL.
func sendActivity(locator Locator, activity mapof.Any) queue.Result {

	const location = "sender.sendActivity"

	// Locate the Actor that is sending this activity
	actorID := activity.GetString(vocab.PropertyActor)
	actor, err := locator.Actor(actorID)

	if err != nil {
		return queue.Failure(derp.Wrap(err, location, "Unable to locate actor", activity))
	}

	// Use the Locator to resolve recipient URIs into inbox URLs
	recipients, err := getRecipients(locator, activity)

	if err != nil {
		return queue.Error(derp.Wrap(err, location, "Unable to retrieve recipient addresses"))
	}

	// Remove duplicate recipient inbox URLs
	uniquer := streams.NewUniquer[string]()
	recipients = uniquer.Range(recipients)

	// Strip BCC and BTo fields before sending
	activity.Remove(vocab.PropertyBCC)
	activity.Remove(vocab.PropertyBTo)

	// Enqueue additional tasks to send this Activity to each recipient's inboxURL
	for inboxURL := range recipients {

		queue.NewTask(OutboxSendActivity, mapof.Any{
			"actorID":  actor.ActorID(),
			"inboxURL": inboxURL,
			"activity": activity,
		})
	}

	// Task Successed Successfully!
	return queue.Success()
}

// sendActivity_SingleRecipient sends a single ActivityPub activity from a the provided Actor to a single recipient's inbox URL.
func sendActivity_SingleRecipient(locator Locator, args mapof.Any) queue.Result {

	const location = "sender.sendActivity"

	// Collect arguments
	actorID := convert.String(args["actorID"])
	inboxURL := convert.String(args["inboxURL"])
	activity := convert.MapOfAny(args["activity"])

	// Locate the Actor that is sending this activity
	actor, err := locator.Actor(actorID)

	if err != nil {
		return queue.Failure(derp.Wrap(err, location, "Unable to retrieve actor for outbound activity", "actorID: "+actorID))
	}

	// Prepare a transaction to send to target Actor's inbox
	transaction := remote.Post(inboxURL).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		With(SignRequest(actor.PrivateKey())).
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
