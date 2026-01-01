package sender

import "crypto"

type testActor struct{}

func (t testActor) ActorID() string {
	return "https://test.actor.social"
}

func (t testActor) PrivateKey() (string, *crypto.PrivateKey) {
	return "https://test.actor.social/#main-key", nil
}
