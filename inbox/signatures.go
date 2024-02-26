package inbox

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
)

/******************************************
 * HTTP Signatures
 *
 * https://docs.joinmastodon.org/spec/security/
 *
 ******************************************/

// validateRequest uses the hannibal/sigs library to verify that the HTTP
// request is signed with a valid key.
func validateRequest(request *http.Request, document streams.Document) error {

	const location = "hannibal.pub.validateRequest"

	// Find the public key for the Actor who signed this request
	keyFinder := keyFinder(document)

	// Verify the request using the Actor's public key
	if err := sigs.Verify(request, keyFinder); err != nil {
		return derp.Wrap(err, location, "Unable to verify HTTP signature", document.Value(), derp.WithCode(derp.CodeForbiddenError))
	}

	return nil
}

// keyFinder looks up the public Key for the provided document/Actor
func keyFinder(document streams.Document) sigs.PublicKeyFinder {

	const location = "hannibal.pub.keyFinder"

	return func(keyID string) (string, error) {

		// Load the Actor from the document
		actor, err := document.Actor().Load()

		if err != nil {
			return "", derp.Wrap(err, location, "Error retrieving Actor from ActivityPub document", document.Value())
		}

		// Search the Actor's public keys for the one that matches the provided keyID
		for key := actor.PublicKey(); key.NotNil(); key = key.Tail() {

			if key.ID() == keyID {
				log.Trace().Str("keyId", keyID).Msg("Hannibal Inbox: Found Public Key")
				return key.PublicKeyPEM(), nil
			}
		}

		// If none match, then return a (hopefully informative) error.
		return "", derp.NewBadRequestError(location, "Actor must publish the key used to sign this request", actor.ID(), keyID)
	}
}
