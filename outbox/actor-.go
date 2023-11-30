package outbox

import (
	"crypto"

	"github.com/benpate/hannibal/queue"
	"github.com/benpate/hannibal/streams"
)

// package-level variable is a singleton queue that is used by default
var defaultQueue queue.Queue

// Actor represents an ActivityPub actor that can send ActivityPub messages
// https://www.w3.org/TR/activitypub/#actors
type Actor struct {

	// Required values passed to NewActor function
	actorID    string
	privateKey crypto.PrivateKey

	// Optional values set via With() options
	publicKeyID string
	followers   <-chan string
	client      streams.Client
	queue       queue.Queue
}

/******************************************
 * Lifecycle Methods
 ******************************************/

// NewActor returns a fully initialized Actor object, and applies optional settings as provided
func NewActor(actorID string, privateKey crypto.PrivateKey, options ...ActorOption) Actor {

	// Set Default Values
	result := Actor{
		actorID:     actorID,
		publicKeyID: actorID + "#main-key",
		privateKey:  privateKey,
	}

	// Apply additional options
	result.With(options...)
	return result
}

// With applies one or more options to an Actor
func (actor *Actor) With(options ...ActorOption) {
	for _, option := range options {
		option(actor)
	}
}

/******************************************
 * Internal / Helper Methods
 ******************************************/

// getQueue returns the queue to use when sending messages
// If the Actor does not include a custom queue, then a default, package-level
// queue is used instead.
func (actor *Actor) getQueue() queue.Queue {

	if actor.queue != nil {
		return actor.queue
	}

	if defaultQueue == nil {
		defaultQueue = queue.NewSimpleQueue(16, 1024)
	}

	return defaultQueue
}

// getClient returns the hannibal Client to use when retrieving
// JSON-LD data.  If the Actor does not include a custom client,
// then a default HTTP-only client is used instead.
func (actor *Actor) getClient() streams.Client {

	if actor.client != nil {
		return actor.client
	}

	return streams.NewDefaultClient()
}
