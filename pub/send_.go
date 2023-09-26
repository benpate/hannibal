package pub

import (
	"encoding/json"
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/mapof"
)

func SendQueueTask(actor Actor, document mapof.Any, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return Send(actor, document, recipient)
	})
}

// Send sends an ActivityStream to a remote ActivityPub service
// actor: The Actor that is sending the request
// message: The ActivityStream that is being sent
// recipient: The remote Actor who will receive the request
func Send(actor Actor, message mapof.Any, recipient streams.Document) error {

	const location = "hannibal.pub.Send"

	// Try to get the inbox from the recipient' profile
	inbox := recipient.Inbox()

	if inbox.IsNil() {
		return derp.NewInternalError(location, "Inbox is empty", recipient.Value())
	}

	inboxID := inbox.ID()

	// Add a "to" field in the message
	message[vocab.PropertyTo] = recipient.ID()

	// Optional debugging output
	if packageDebugLevel >= DebugLevelTerse {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("------------------------------------------")
		}
		fmt.Println("HANNIBAL: Sending Activity: " + recipient.ID())
		if packageDebugLevel >= DebugLevelVerbose {
			marshalled, _ := json.MarshalIndent(message, "", "  ")
			fmt.Println(string(marshalled))
		}
	}

	// Send the request to the target Actor's inbox
	transaction := remote.Post(inboxID).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		Use(RequestSignature(actor)).
		JSON(message)

	if packageDebugLevel >= DebugLevelVerbose {
		transaction.Use(options.Debug())
	}

	if err := transaction.Send(); err != nil {
		return derp.Wrap(err, location, "Error sending Follow request", inbox)
	}

	if packageDebugLevel >= DebugLevelVerbose {
		fmt.Println("HANNIBAL: Sent Activity Successfully")
	}

	// Done!
	return nil
}
