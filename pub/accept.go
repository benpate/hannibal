package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

func SendAcceptQueueTask(actor Actor, activity streams.Document) QueueTask {
	return NewQueueTask(func() error {
		return SendAccept(actor, activity)
	})
}

func SendAccept(actor Actor, activity streams.Document) error {

	message := mapof.Any{
		"@context": vocab.ContextTypeActivityStreams,
		"type":     vocab.ActivityTypeAccept,
		"actor":    actor.ActorID,
		"object":   activity.Value(),
	}

	targetURL := activity.Actor().ID()

	if err := Send(actor, message, targetURL); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
