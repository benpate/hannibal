package sigs

// PublicKeyFinder is a function that can look up a public key.
// This is injected into the Verify function by the inbox.
type PublicKeyFinder func(keyID string) (string, error)
