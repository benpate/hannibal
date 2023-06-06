package streams

import (
	"github.com/benpate/derp"
	"github.com/benpate/remote"
)

type DefaultClient struct{}

func NewDefaultClient() Client {
	return DefaultClient{}
}

func (client DefaultClient) Load(uri string) (Document, error) {

	// Try to load-and-parse the value from the remote server
	result := make(map[string]any)

	transaction := remote.Get(uri).
		Accept("application/activity+json").
		Response(&result, nil)

	if err := transaction.Send(); err != nil {
		return NilDocument(), derp.Wrap(err, "hannibal.streams.Client.Load", "Error loading JSON-LD document", uri)
	}

	header := transaction.ResponseObject.Header

	// Return in triumph
	return NewDocument(result,
			WithClient(client),
			WithMeta("cache-control", header.Get("cache-control")),
			WithMeta("etag", header.Get("etag")),
			WithMeta("expires", header.Get("expires")),
		),
		nil
}
