package pub

import (
	"fmt"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/hannibal/streams"
	"github.com/davecgh/go-spew/spew"
)

/******************************************
 * HTTP Signatures
 *
 * https://docs.joinmastodon.org/spec/security/
 *
 ******************************************/

// validateRequest verifies that the HTTP request is signed with a valid key.
// This function loads the public key from the ActivityPub actor, then verifies their signature.
func validateRequest(request *http.Request, document streams.Document) error {

	// TODO: HIGH: Validate http Signature headers
	// TODO: HIGH: Validate Digest headers
	// TODO: HIGH: Confirm that the http signature includes "(request-target)" "host" "date" and "digest" (extras are ok)

	// Add required "host" header if it doesn't already exist
	request.Header.Set("host", request.Host)

	const location = "hannibal.pub.validateRequest"

	// Get the Actor from the document
	actor, err := document.Actor().Load()

	if err != nil {
		err = derp.Wrap(err, location, "Error retrieving Actor from ActivityPub document")
		derp.Report(err)
		return err
	}

	// Get the Actor's Public Key
	actorPublicKey, err := actor.PublicKey().Load()

	if err != nil {
		return derp.Wrap(err, location, "Error retrieving Public Key from Actor")
	}

	actorPublicPEM := actorPublicKey.PublicKeyPEM()

	if packageDebugLevel >= DebugLevelVerbose {
		fmt.Println("------------------------------------------")
		fmt.Println(location)
		spew.Dump(actor.Value())
		fmt.Println("PEM: " + actorPublicPEM)
	}

	// Verify the request using the Actor's public key
	if err := sigs.Verify(request, actorPublicPEM); err != nil {
		derp.SetErrorCode(err, derp.CodeForbiddenError)
		return derp.Wrap(err, location, "Unable to verify HTTP signature")
	}

	return nil
}
