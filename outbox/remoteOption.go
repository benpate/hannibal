package outbox

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/remote"
)

// SignRequest is a middleware for the remote package that adds an HTTP Signature to a request.
func SignRequest(actor Actor) remote.Option {

	return remote.Option{

		ModifyRequest: func(txn *remote.Transaction, request *http.Request) *http.Response {

			// Add a "Digest" header to the request and sign the outgoing request.
			if err := sigs.Sign(request, actor.publicKeyID, actor.privateKey); err != nil {
				derp.Report(derp.Wrap(err, "activitypub.RequestSignature", "Error signing HTTP request.  This is likely because of a problem with the actor's private key."))
			}

			// If exists, write the Digest back into the transaction (for serialization, et al)
			if digest := request.Header.Get("Digest"); digest != "" {
				txn.Header("Digest", digest)
			}

			// If exists, write the Signature back into the transaction (for serialization, et al)
			if signature := request.Header.Get("Signature"); signature != "" {
				txn.Header("Signature", signature)
			}

			// Oh, yeah...
			return nil
		},
	}
}
