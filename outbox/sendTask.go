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

	logger := log.With().Str("loc", location).Logger()

	logger.Info().Str("recipient", task.recipient.ID()).Msg("Sending activity...")

	if canLog(zerolog.DebugLevel) {
		rawJSON, _ := json.MarshalIndent(task.message, "", "  ")
		logger.Debug().Msg(string(rawJSON))
	}

	inboxURL := task.recipient.Inbox().ID()

	if inboxURL == "" {
		logger.Error().Msg("Recipient does not have an inbox")
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

	logger.Debug().Msg("Activity sent successfully")

	// Done!
	return nil
}
