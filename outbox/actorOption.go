package outbox

import (
	"github.com/benpate/hannibal/streams"
)

// ActorOption is a function signature that modifies optional settings for an Actor
type ActorOption func(*Actor)

// WithPublicKey is an ActorOption that sets the public key for an Actor
func WithPublicKey(publicKeyID string) ActorOption {
	return func(a *Actor) {
		a.publicKeyID = publicKeyID
	}
}

// WithCliient is an ActorOption that sets the hanibal Client for an Actor
func WithClient(client streams.Client) ActorOption {
	return func(a *Actor) {
		a.client = client
	}
}

// TODO: Restore Queue::
/*
// WithQueue is an ActorOption that sets the outbound Queue for an Actor
func WithQueue(queue *queue.Queue) ActorOption {
	return func(a *Actor) {
		a.queue = queue
	}
}
*/

// WithFollowers is an ActorOption that provides a channel of followers for an Actor
func WithFollowers(followers <-chan string) ActorOption {
	return func(a *Actor) {
		a.followers = followers
	}
}
