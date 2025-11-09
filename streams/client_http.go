package streams

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
)

type DefaultClient struct{}

func NewDefaultClient() Client {
	return DefaultClient{}
}

func (client DefaultClient) SetRootClient(rootClient Client) {}

// Load implements the hannibal.Client interface, which loads an ActivityStream
// document from a remote server. For the hannibal default client, this method
// simply loads the document from a remote server with no other processing.
func (client DefaultClient) Load(url string, options ...any) (Document, error) {

	const location = "hannibal.streams.Client.Load"

	result := make(map[string]any)

	// Try to load-and-parse the value from the remote server
	transaction := remote.Get(url).
		Accept(vocab.ContentTypeActivityPub).
		Result(&result)

	if err := transaction.Send(); err != nil {
		return NilDocument(), derp.Wrap(err, location, "Unable to load JSON-LD document", url)
	}

	// Return in triumph
	return NewDocument(result,
			WithClient(client),
			WithHTTPHeader(transaction.ResponseHeader()),
		),
		nil
}

// Save is required to implement the document.Cache interface.
// For this client, Save is a NOOP
func (client DefaultClient) Save(document Document) error {
	return nil
}

// Delete is required to implement the document.Cache interface.
// For this client, Delete is a NOOP
func (client DefaultClient) Delete(documentID string) error {
	return nil
}
