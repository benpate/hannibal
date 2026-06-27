package outbox

import (
	"iter"

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

// WithClient is an ActorOption that sets the Hannibal Client for an Actor.
func WithClient(client streams.Client) ActorOption {
	return func(a *Actor) {
		a.client = client
	}
}

// WithFollowers is an ActorOption that sets the Actor's followers iterator.
func WithFollowers(followers iter.Seq[string]) ActorOption {
	return func(a *Actor) {
		a.followers = followers
	}
}
