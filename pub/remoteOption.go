package pub

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/remote"
)

// RequestSiguature is a middleware for the remote package that adds an HTTP Signature to a request.
func RequestSignature(actor Actor) remote.Option {

	return remote.Option{

		ModifyRequest: func(_ *remote.Transaction, request *http.Request) *http.Response {

			// Add a "Digest" header to the request and sign the outgoing request.
			if err := sigs.Sign(request, actor.PublicKeyID, actor.PrivateKey); err != nil {
				derp.Report(derp.Wrap(err, "activitypub.RequestSignature", "Error signing HTTP request.  This is likely because of a problem with the actor's private key."))
			}

			// Oh, yeah...
			return nil
		},
	}
}
