package pub

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/remote"
	"github.com/go-fed/httpsig"
)

func RequestMakeDigest() remote.Middleware {

	return remote.Middleware{

		Request: func(request *http.Request) error {
			var body bytes.Buffer
			bodyReader, err := request.GetBody()

			if err != nil {
				return derp.Wrap(err, "pub.RequestDigest", "Error getting request body")
			}

			if _, err := body.ReadFrom(bodyReader); err != nil {
				return derp.Wrap(err, "pub.RequestDigest", "Error reading request body")
			}

			digest := sha256.Sum256(body.Bytes())
			request.Header.Set("digest", "SHA-256="+base64.StdEncoding.EncodeToString(digest[:]))

			return nil
		},
	}
}

// RequestSiguature is a middleware for the remote package that adds an HTTP Signature to a request.
func RequestSignature(actor Actor) remote.Middleware {

	return remote.Middleware{

		Config: func(t *remote.Transaction) error {
			return nil
		},

		Request: func(request *http.Request) error {

			// Add required headers if they don't already exist
			request.Header.Set("host", request.URL.Host)
			request.Header.Set("date", time.Now().UTC().Format(http.TimeFormat))

			// Collect settings to sign the request
			preferredAlgorithms := []httpsig.Algorithm{httpsig.RSA_SHA256, httpsig.RSA_SHA512}
			defaultAlgorithm := httpsig.DigestAlgorithm(httpsig.RSA_SHA512)
			headers := []string{"(request-target)", "host", "date", "digest"}

			// Try to make a new signer
			signer, _, err := httpsig.NewSigner(preferredAlgorithms, defaultAlgorithm, headers, httpsig.Signature, 30)

			if err != nil {
				return derp.Wrap(err, "activitypub.RequestSignature", "Error creating HTTP Signature signer")
			}

			// Sign the request
			if err := signer.SignRequest(actor.PrivateKey, actor.PublicKeyID, request, nil); err != nil {
				return derp.Wrap(err, "activitypub.RequestSignature", "Error signing HTTP request")
			}

			// Oh, yeah...
			return nil
		},

		Response: func(r *http.Response, b *[]byte) error {
			return nil
		},
	}
}
