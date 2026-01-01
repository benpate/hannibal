package sender

import "crypto"

// Actor defines the interface for an ActivityPub Actor
// Actors must be addressable
// https://www.w3.org/TR/activitypub/#actors
type Actor interface {

	// ActorID returns the unique ID (URL) of this Actor
	ActorID() string

	// PrivateKey returns a PrivateKeyID and PrivateKey to use for
	// signing outbound ActivityPub messages from this Actor
	PrivateKey() (privateKeyID string, privateKey *crypto.PrivateKey)
}

func NewActor(actorID string, publicKeyID string, privateKey crypto.PrivateKey) Actor {
	return defaultActor{
		actorID:     actorID,
		publicKeyID: publicKeyID,
		privateKey:  privateKey,
	}
}

// defaultActor is a basic implementation of the Actor interface
type defaultActor struct {
	actorID     string
	publicKeyID string
	privateKey  crypto.PrivateKey
}

// ActorID is a part of the sender.Actor interface
// It returns the unique ID (URL) of this Actor
func (actor defaultActor) ActorID() string {
	return actor.actorID
}

// PrivateKey is a part of the sender.Actor interface
// It returns a PrivateKeyID and PrivateKey to use for signing outbound ActivityPub messages from this Actor
func (actor defaultActor) PrivateKey() (publicKeyID string, privateKey *crypto.PrivateKey) {
	return actor.publicKeyID, &actor.privateKey
}
