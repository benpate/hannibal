package sigs

import (
	"net/http"
)

// GetAuthenticatedActor returns the Actor ID from a verified HTTP Signature.
func GetAuthenticatedActor(r *http.Request, publicKeyFinder PublicKeyFinder) string {

	// If we have a valid HTTP signature, then use it for the authenticatedID
	if signature, err := Verify(r, publicKeyFinder); err == nil {
		return signature.ActorID()
	}

	return ""
}
