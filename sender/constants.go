package sender

// OutboxSendActivity is the name of the task that scans an
// outbound ActivityPub activity for recipients' inbox URLs,
// then queues additional tasks to send that activity to
// each recipient
const OutboxSendActivity = "Outbox:SendActivity"

// OutboxSendActivity_SingleRecipient is the name of the task that sends
// an outbound ActivityPub activity to a single recipient.
const OutboxSendActivity_SingleRecipient = "Outbox:SendActivity:SingleRecipient"
