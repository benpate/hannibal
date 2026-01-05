package sender

import (
	"github.com/benpate/turbine/queue"
)

// Consumer returns a turbine queue.Consumer that processes
// outbound ActivityPub actitities for this outbox.
func Consumer(sender Sender) queue.Consumer {

	return func(name string, args map[string]any) queue.Result {

		switch name {

		// Catalog all recipients and queue individual send tasks
		case OutboxSendToAllRecipients:
			return sender.SendToAllRecipients(args)

		// Send an activity to a single recipient
		case OutboxSendToSingleRecipient:
			return sender.SendToSingleRecipient(args)
		}

		// All other task names are left for other consumers.
		return queue.Ignored()
	}
}
