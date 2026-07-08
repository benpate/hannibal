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

// WithAllowPrivateIPs is an ActorOption that permits (or forbids) outbound
// delivery to non-public (private/loopback) addresses. It is FALSE by default so
// that remote's SSRF guard stays active; callers delivering to a local peer (for
// example, a dev instance federating with itself) opt in by passing TRUE.
func WithAllowPrivateIPs(allowPrivateIPs bool) ActorOption {
	return func(a *Actor) {
		a.allowPrivateIPs = allowPrivateIPs
	}
}
