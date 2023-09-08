package pub

import (
	"io"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/remote"
)

// RequestSiguature is a middleware for the remote package that adds an HTTP Signature to a request.
func RequestSignature(actor Actor) remote.Middleware {

	return remote.Middleware{

		Config: func(t *remote.Transaction) error {
			return nil
		},

		Request: func(request *http.Request) error {

			bodyReader, err := request.GetBody()

			if err != nil {
				return derp.Wrap(err, "activitypub.RequestSignature", "Error getting body from request")
			}

			body, err := io.ReadAll(bodyReader)

			if err != nil {
				return derp.Wrap(err, "activitypub.RequestSignature", "Error reading body from request")
			}

			// Sign the outgoing request.  This also adds a "Digest" header to the request.
			if err := sigs.Sign(request, body, actor.PublicKeyID, actor.PrivateKey); err != nil {
				return derp.Wrap(err, "activitypub.RequestSignature", "Error signing HTTP request.  This is likely because of a problem with the actor's private key.")
			}

			// Oh, yeah...
			return nil
		},

		Response: func(r *http.Response, b *[]byte) error {
			return nil
		},
	}
}
