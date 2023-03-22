package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

func SendUpdateQueueTask(actor Actor, activity mapof.Any, targetURL string) QueueTask {
	return NewQueueTask(func() error {
		return SendUpdate(actor, activity, targetURL)
	})
}

func SendUpdate(actor Actor, activity mapof.Any, targetURL string) error {

	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeUpdate,
		"actor":    actor.ActorID,
		"object":   activity,
	}

	if err := Send(actor, message, targetURL); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
