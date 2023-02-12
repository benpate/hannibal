package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/jsonld"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

func PostAcceptQueueTask(actor Actor, activity jsonld.Reader) QueueTask {
	return NewQueueTask(func() error {
		return PostAccept(actor, activity)
	})
}

func PostAccept(actor Actor, activity jsonld.Reader) error {

	message := mapof.Any{
		"@context": DefaultContext,
		"type":     vocab.ActivityTypeAccept,
		"actor":    actor.ActorID,
		"object":   activity,
	}

	targetURL := activity.ActorID()

	if err := Post(actor, message, targetURL); err != nil {
		return derp.Wrap(err, "activitypub.PostAcceptActivity", "Error sending Accept request")
	}

	return nil
}
