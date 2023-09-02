package pub

import (
	"crypto/rand"
	"crypto/rsa"
)

func getTestActor() Actor {

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return Actor{
		ActorID:     "https://example.com/actor/1",
		PublicKeyID: "https://example.com/actor/1#main-key",
		PublicKey:   privateKey.PublicKey,
		PrivateKey:  privateKey,
	}
}
