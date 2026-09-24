package clients

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// CarpoolClient is a streams.Client wrapper that shares each in-progress Load with every other
// caller in its Carpool that loads the same URL as the same signer.
type CarpoolClient struct {
	carpool     *Carpool
	innerClient streams.Client
	rootClient  streams.Client
	signer      string
}

// Load retrieves a document from the innerClient, joining a Load of the same URL by the same signer
// that is already in progress anywhere in the Carpool.
func (client *CarpoolClient) Load(url string, options ...any) (streams.Document, error) {

	const location = "hannibal.clients.CarpoolClient.Load"

	// RULE: Options are opaque here, so two calls with options cannot be proven to want the same
	// result.  A caller that passes any option (a forced refresh, a reveal) always loads alone.
	if len(options) > 0 {
		return client.innerClient.Load(url, options...)
	}

	// RULE: A remote server may answer each signer differently, so riders share a Load only when
	// they would have signed it the same way.  Signers never contain NUL, so the key is unambiguous.
	key := client.signer + "\x00" + url

	// Load the document, or wait for the caller that is already loading it
	result, err, shared := client.carpool.group.Do(key, func() (any, error) {
		return client.innerClient.Load(url)
	})

	document, isDocument := result.(streams.Document)

	if !isDocument {
		document = streams.NilDocument()
	}

	// RULE: A shared result is never handed out directly.  Each rider gets a copy bound to its
	// own client stack, so no two callers share maps or sign follow-up loads as each other.
	if shared {
		document = client.rebind(document)
	}

	if err != nil {
		return document, derp.Wrap(err, location, "Loading document", url)
	}

	// Arrived safely
	return document, nil
}

// Save stores the Document in the underlying cache.
func (client *CarpoolClient) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

// Delete removes a document from the underlying client's cache.
func (client *CarpoolClient) Delete(documentID string) error {
	return client.innerClient.Delete(documentID)
}

// SetRootClient records the top-level client and passes it down to the underlying client.
func (client *CarpoolClient) SetRootClient(rootClient streams.Client) {

	client.rootClient = rootClient

	if client.innerClient != nil {
		client.innerClient.SetRootClient(rootClient)
	}
}

// rebind returns a deep copy of a shared document, bound to this caller's own client stack
func (client *CarpoolClient) rebind(document streams.Document) streams.Document {

	// A stack with no root reported yet resolves follow-up loads through this client
	var rootClient streams.Client = client

	if client.rootClient != nil {
		rootClient = client.rootClient
	}

	return document.Clone().AddOptions(streams.WithClient(rootClient))
}

// Verify that CarpoolClient satisfies the streams.Client interface.
var _ streams.Client = &CarpoolClient{}
