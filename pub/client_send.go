package pub

import (
	"encoding/json"
	"fmt"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

func SendQueueTask(actor Actor, document mapof.Any, targetID string) QueueTask {
	return NewQueueTask(func() error {
		return Send(actor, document, targetID)
	})
}

// Send sends an ActivityStream to a remote ActivityPub service
// actor: The Actor that is sending the request
// document: The ActivityStream that is being sent
// targetID: The ID of the Actor that will receive the request
func Send(actor Actor, document mapof.Any, targetID string) error {

	const location = "hannibal.pub.Send"

	// Optional debugging output
	if packageDebugLevel >= DebugLevelTerse {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("------------------------------------------")
		}
		fmt.Println("HANNIBAL: Sending Activity: " + targetID)
		if packageDebugLevel >= DebugLevelVerbose {
			marshalled, _ := json.MarshalIndent(document, "", "  ")
			fmt.Println(string(marshalled))
		}
	}

	// Try to get the source profile that we're going to follow
	target, err := GetProfile(targetID)

	if err != nil {
		return derp.Wrap(err, location, "Error getting source profile", targetID)
	}

	// Try to get the actor's inbox from the actor ActivityStream.
	// TODO: LOW: Is there a better / more reliable way to do this?
	inbox := target.GetString("inbox")

	if inbox == "" {
		return derp.NewInternalError(location, "Unable to find 'inbox' in target profile", targetID, target)
	}

	// Send the request to the target Actor's inbox
	transaction := remote.Post(inbox).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		Use(RequestSignature(actor)).
		JSON(document)

	if err := transaction.Send(); err != nil {
		return derp.Wrap(err, location, "Error sending Follow request", inbox)
	}

	if packageDebugLevel >= DebugLevelVerbose {
		fmt.Println("HANNIBAL: Sent Acivity Successfully")
	}

	// Done!
	return nil
}

func SendActivityQueueTask(actor Actor, activityType string, object any, targetID string) QueueTask {
	return NewQueueTask(func() error {
		return SendActivity(actor, activityType, object, targetID)
	})
}

// SendActivity wraps a document in a standard ActivityStream envelope and sends it to the target.
func SendActivity(actor Actor, activityType string, object any, targetID string) error {

	document := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     activityType,
		"actor":    actor.ActorID,
		"object":   object,
	}

	return Send(actor, document, targetID)
}
