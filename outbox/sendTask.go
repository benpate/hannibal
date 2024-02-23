package outbox

import (
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SendTask implements the queue.Task interface
type SendTask struct {
	actor     Actor
	message   mapof.Any
	recipient streams.Document
}

// NewSendTask returns a fully initialized SendTask object
func NewSendTask(actor Actor, message mapof.Any, recipient streams.Document) SendTask {
	return SendTask{
		actor:     actor,
		message:   message,
		recipient: recipient,
	}
}

func (task SendTask) Run() error {

	const location = "hannibal.outbox.SendTask.Run"

	inboxURL := task.recipient.Inbox().ID()

	if canLog(zerolog.DebugLevel) {
		messageID := task.message.GetString(vocab.PropertyID)
		log.Debug().Str("loc", location).Str("id", messageID).Str("to", inboxURL).Msg("Sending:")

		if canLog(zerolog.TraceLevel) {
			rawJSON, _ := json.MarshalIndent(task.message, "", "  ")
			log.Trace().Msg(string(rawJSON))
		}
	}

	if inboxURL == "" {
		log.Error().Msg("Recipient does not have an inbox")
		return nil // returning nil error because we have failed so bacly that we don't even want to retry.
	}

	// Send the request to the target Actor's inbox
	transaction := remote.Post(inboxURL).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		Use(SignRequest(task.actor)).
		JSON(task.message)

	if canLog(zerolog.TraceLevel) {
		transaction.Use(options.Debug())
	}

	if err := transaction.Send(); err != nil {
		return derp.ReportAndReturn(derp.Wrap(err, location, "Error sending ActivityPub request", inboxURL))
	}

	log.Debug().Msg("Outbox: sendTask: Activity sent successfully")

	// Done!
	return nil
}
