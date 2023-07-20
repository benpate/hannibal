package streams

import (
	"github.com/benpate/derp"
	"github.com/benpate/remote"
)

type DefaultClient struct{}

func NewDefaultClient() Client {
	return DefaultClient{}
}

func (client DefaultClient) LoadActor(uri string) (Document, error) {
	return client.LoadDocument(uri, map[string]any{})
}

func (client DefaultClient) LoadDocument(uri string, defaultValue map[string]any) (Document, error) {

	// Try to load-and-parse the value from the remote server
	transaction := remote.Get(uri).
		Accept("application/activity+json").
		Response(&defaultValue, nil)

	if err := transaction.Send(); err != nil {
		return NilDocument(), derp.Wrap(err, "hannibal.streams.Client.Load", "Error loading JSON-LD document", uri)
	}

	header := transaction.ResponseObject.Header

	// Return in triumph
	return NewDocument(defaultValue,
			WithClient(client),
			WithMeta("cache-control", header.Get("cache-control")),
			WithMeta("etag", header.Get("etag")),
			WithMeta("expires", header.Get("expires")),
		),
		nil
}
