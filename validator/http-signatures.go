package validator

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
)

// HTTPSig is a Validator that checks incoming HTTP requests
// using the HTTP signatures algorithm.
// https://docs.joinmastodon.org/spec/security/
type HTTPSig struct{}

func NewHTTPSig() HTTPSig {
	return HTTPSig{}
}

// Validate uses the hannibal/sigs library to verify that the HTTP
// request is signed with a valid key.
func (validator HTTPSig) Validate(request *http.Request, document *streams.Document) Result {

	if !sigs.HasSignature(request) {
		return ResultUnknown
	}

	// Find the public key for the Actor who signed this request
	keyFinder := validator.keyFinder(document)

	// Verify the request using the Actor's public key
	if err := sigs.Verify(request, keyFinder); err != nil {
		log.Trace().Err(err).Msg("Hannibal Inbox: Error verifying HTTP Signature")
		return ResultInvalid
	}

	log.Trace().Msg("Hannibal Inbox: HTTP Signature Verified")
	return ResultValid
}

// keyFinder looks up the public Key for the provided document/Actor using the
// HTTP client in the document.
func (validator HTTPSig) keyFinder(document *streams.Document) sigs.PublicKeyFinder {

	const location = "hannibal.validator.HTTPSig.keyFinder"

	return func(keyID string) (string, error) {

		// Load the Actor from the document
		actor, err := document.Actor().Load()

		if err != nil {
			return "", derp.Wrap(err, location, "Error retrieving Actor from ActivityPub document", document.Value())
		}

		// Search the Actor's public keys for the one that matches the provided keyID
		for key := actor.PublicKey(); key.NotNil(); key = key.Tail() {

			if key.ID() == keyID {
				return key.PublicKeyPEM(), nil
			}
		}

		// If none match, then return a (hopefully informative) error.
		log.Trace().Str("keyId", keyID).Msg("Hannibal Inbox: Could not find remote actor's public key")
		return "", derp.NewBadRequestError(location, "Actor must publish the key used to sign this request", actor.ID(), keyID)
	}
}
