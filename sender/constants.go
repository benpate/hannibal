package sender

// OutboxSendToAllRecipients is the name of the task that scans an
// outbound ActivityPub activity for recipients' inbox URLs,
// then queues additional tasks to send that activity to
// each recipient
const OutboxSendToAllRecipients = "Outbox:SendToAllRecipients"

// OutboxSendToSingleRecipient is the name of the task that sends
// an outbound ActivityPub activity to a single recipient.
const OutboxSendToSingleRecipient = "Outbox:SendToSingleRecipient"
