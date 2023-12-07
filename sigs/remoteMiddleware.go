package sigs

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/remote"
)

// WithSigner is a remote.Option that signs an outbound HTTP request
func WithSigner(signer Signer) remote.Option {

	return remote.Option{

		ModifyRequest: func(txn *remote.Transaction, request *http.Request) *http.Response {

			// Sign the outbound request
			if err := signer.Sign(request); err != nil {
				derp.Report(derp.Wrap(err, "hannibal.sigs.WithSigner", "Error signing request"))
			}

			// Nil response means that we are still sending the request to the remote server
			// instead of replacing it with a new request.
			return nil
		},
	}
}
