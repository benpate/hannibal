package pub

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/davecgh/go-spew/spew"
)

// ReceiveInboxRequest reads an incoming HTTP request and returns a parsed and validated ActivityPub activity
func ReceiveInboxRequest(request *http.Request, cache streams.Cache) (document streams.Document, err error) {

	const location = "activitypub.ReceiveInboxRequest"

	// RULE: Content-Type MUST be "application/activity+json" or "application/ld+json"
	if !IsActivityPubContentType(request.Header.Get(vocab.ContentType)) {
		return streams.NilDocument(), derp.NewBadRequestError(location, "Content-Type MUST be 'application/activity+json'")
	}

	spew.Dump("Received ActivityPub request", request.Header)

	// Try to read the body from the request
	var bodyBuffer bytes.Buffer
	if _, err := bodyBuffer.ReadFrom(request.Body); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error reading body into buffer")
	}

	spew.Dump(bodyBuffer.String())

	// Try to retrieve the object from the buffer
	document = streams.NewDocument(nil, cache)

	if err := json.Unmarshal(bodyBuffer.Bytes(), &document); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error unmarshalling JSON body into ActivityPub document")
	}

	// Validate the Actor and Public Key
	if err := validateRequest(request, document, &bodyBuffer, request.Header.Get("Signature")); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Request is invalid")
	}

	documentType := document.Type()

	// First, assume that we have a fully defined activity
	if activityType := vocab.ValidateActivityType(documentType); activityType != vocab.Unknown {
		return document, nil
	}

	// Otherwise, assume that we have an implicit "Create" activity
	if objectType := vocab.ValidateObjectType(documentType); objectType != vocab.Unknown {
		// TODO: MEDIUM: Wrap original activity in a "Create" activity
		// TODO: MEDIUM: Can we get the Actor from the signed HTTP Request?
		return document, derp.NewInternalError(location, "Implicit 'Create' activities are not yet implemented")
	}

	// Return the activity to the caller.
	return streams.NilDocument(), derp.NewInternalError(location, "Unknown ActivityPub message type", document.Value())
}

func validateRequest(request *http.Request, document streams.Document, bodyBuffer *bytes.Buffer, httpHeaderSignature string) error {
	/*

		const location = "activitypub.validateRequest"

		signature := signatures.ParseSignatureHeader(httpHeaderSignature)

		// Get the Actor from the document
		actor, err := document.Actor().AsObject()

		if err != nil {
			return derp.Wrap(err, location, "Error retrieving Actor from ActivityPub document")
		}

		// Get the Actor's Public Key
		actorPublicKey, err := actor.PublicKey().AsObject()

		if err != nil {
			return derp.Wrap(err, location, "Error retrieving Public Key from Actor")
		}

		actorPublicKeyPEM := actorPublicKey.PublicKeyPEM()

		// Parse the Public Key
		key, err := ssh.ParsePublicKey([]byte(actorPublicKeyPEM))

		if err != nil {
			return derp.Wrap(err, location, "Error parsing Public Key")
		}

		// Finally, Verify request signatures
		verifier, err := httpsig.NewVerifier(request)

		if err != nil {
			return derp.Wrap(err, location, "Error creating HTTP Signature verifier")
		}

		algorithm := httpsig.Algorithm(signature.GetString("algorithm"))

		if err := verifier.Verify(key, algorithm); err != nil {
			return derp.Wrap(err, location, "Error verifying HTTP Signature")
		}
	*/
	return nil
}
