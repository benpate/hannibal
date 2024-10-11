package outbox

import (
	"crypto"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/turbine/queue"
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
	queue       *queue.Queue
}

/******************************************
 * Lifecycle Methods
 ******************************************/

// NewActor returns a fully initialized Actor object, and applies optional settings as provided
func NewActor(actorID string, privateKey crypto.PrivateKey, queue *queue.Queue, options ...ActorOption) Actor {

	// Set Default Values
	result := Actor{
		actorID:     actorID,
		publicKeyID: actorID + "#main-key",
		privateKey:  privateKey,
		queue:       queue,
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

func (actor *Actor) ActorID() string {
	return actor.actorID
}

/******************************************
 * Internal / Helper Methods
 ******************************************/

// getClient returns the hannibal Client to use when retrieving
// JSON-LD data.  If the Actor does not include a custom client,
// then a default HTTP-only client is used instead.
func (actor *Actor) getClient() streams.Client {

	if actor.client != nil {
		return actor.client
	}

	return streams.NewDefaultClient()
}
