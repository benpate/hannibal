package pub

import (
	"encoding/json"
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/middleware"
	"github.com/benpate/rosetta/mapof"
)

func SendQueueTask(actor Actor, document mapof.Any, recipient streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return Send(actor, document, recipient)
	})
}

// Send sends an ActivityStream to a remote ActivityPub service
// actor: The Actor that is sending the request
// document: The ActivityStream that is being sent
// recipient: The remote Actor who will receive the request
func Send(actor Actor, document mapof.Any, recipient streams.Document) error {

	const location = "hannibal.pub.Send"

	// Try to get the inbox from the recipient' profile
	inbox := recipient.Inbox().ID()

	if inbox == "" {
		return derp.NewInternalError(location, "Inbox is empty", recipient.Value())
	}

	// Optional debugging output
	if packageDebugLevel >= DebugLevelTerse {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("------------------------------------------")
		}
		fmt.Println("HANNIBAL: Sending Activity: " + recipient.ID())
		if packageDebugLevel >= DebugLevelVerbose {
			marshalled, _ := json.MarshalIndent(document, "", "  ")
			fmt.Println(string(marshalled))
		}
	}

	// Add a "to" field in the document
	document[vocab.PropertyTo] = inbox

	/*/ Marshal the document into a byte reader
	bodyBytes, err := json.Marshal(document)

	spew.Dump(string(bodyBytes), len(bodyBytes), err)
	if err != nil {
		return derp.Wrap(err, location, "Error marshalling ActivityStream", document)
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", inbox, bytes.NewReader(bodyBytes))

	if err != nil {
		return derp.Wrap(err, location, "Error creating HTTP Request", inbox)
	}

	req.Header.Set("Accept", vocab.ContentTypeActivityPub)
	req.Header.Set("Content-Type", vocab.ContentTypeActivityPub)

	middleware := RequestSignature(actor)
	if err := middleware.Request(req); err != nil {
		return derp.Wrap(err, location, "Error signing request", req)
	}

	spew.Dump("HERE?? ==========")


	client := http.Client{}
	spew.Dump("A")
	response, err := client.Do(req)
	spew.Dump("B")

	if err != nil {
		spew.Dump("C")
		spew.Dump("FAIL", err)
		return derp.Wrap(err, location, "Error sending HTTP Request", req)
	}

	spew.Dump("D")
	spew.Dump("SUCCESS", response.Header)

	responseBytes, err := io.ReadAll(response.Body)

	if err != nil {
		spew.Dump("E", err)
		return derp.Wrap(err, location, "Error reading HTTP Response", response)
	}

	spew.Dump(string(responseBytes))
	*/

	// Send the request to the target Actor's inbox
	transaction := remote.Post(inbox).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		Use(RequestSignature(actor)).
		JSON(document)

	if packageDebugLevel >= DebugLevelVerbose {
		transaction.Use(middleware.Debug())
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
