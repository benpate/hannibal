package inbox

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

// validateRequest uses the hannibal/sigs library to verify that the HTTP
// request is signed with a valid key.
func validateRequest(request *http.Request, document streams.Document) error {

	const location = "hannibal.pub.validateRequest"

	// Get the Actor from the document
	actor, err := document.Actor().Load()

	if err != nil {
		err = derp.Wrap(err, location, "Error retrieving Actor from ActivityPub document", document.Value())
		derp.Report(err)
		return err
	}

	// Get the Actor's Public Key
	actorPublicKey := actor.PublicKey().LoadLink()
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
		return derp.Wrap(err, location, "Unable to verify HTTP signature", document.Value())
	}

	return nil
}
