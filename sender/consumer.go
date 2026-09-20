package sender

import (
	"github.com/benpate/turbine/queue"
)

// outboxConsumer is the turbine queue.Consumer that processes outbound ActivityPub
// activities for one outbox.
type outboxConsumer struct {
	sender Sender
}

// Consumer returns a turbine queue.Consumer that processes
// outbound ActivityPub activities for this outbox.
func Consumer(sender Sender) queue.Consumer {
	return &outboxConsumer{
		sender: sender,
	}
}

// OnPublish implements the queue.Consumer interface.
func (consumer *outboxConsumer) OnPublish(*queue.Task) error {
	// An outbound activity needs nothing changed on its way onto the queue.
	return nil
}

// Run executes a single attempt of an outbound send task.
// Implements the queue.Consumer interface.
func (consumer *outboxConsumer) Run(task queue.Task) queue.Result {

	switch task.Name {

	// Catalog all recipients and queue individual send tasks
	case OutboxSendToAllRecipients:
		return consumer.sender.SendToAllRecipients(task.Arguments)

	// Send an activity to a single recipient
	case OutboxSendToSingleRecipient:
		return consumer.sender.SendToSingleRecipient(task.Arguments)
	}

	// All other task names are left for other consumers.
	return queue.Ignored()
}

// OnSuccess implements the queue.Consumer interface.
func (consumer *outboxConsumer) OnSuccess(queue.Task) error {
	// A delivered activity reports nothing.
	return nil
}

// OnError implements the queue.Consumer interface.
func (consumer *outboxConsumer) OnError(queue.Task, error) error {
	// A delivery that will be retried reports nothing; SendToSingleRecipient owns the retry policy.
	return nil
}

// OnFailure implements the queue.Consumer interface.
func (consumer *outboxConsumer) OnFailure(queue.Task, error) error {
	// An undeliverable activity reports nothing yet.
	return nil
}
