package jsonld

import (
	"github.com/benpate/derp"
	"github.com/benpate/remote"
)

// HTTPClient wraps http transactions to load remote JSON-LD documents.
type HTTPClient struct {
	mimeType string
	cache    Cache
}

// NewClient creates a new HTTPClient object, which can be used to load remote JSON-LD documents.
func NewClient(mimeType string, cache Cache) HTTPClient {
	return HTTPClient{
		mimeType: mimeType,
		cache:    cache,
	}
}

// Load retrieves a JSON-LD document from a remote server, parses is, and returns a Reader object.
func (client *HTTPClient) Load(uri string) (Reader, error) {

	// If the value exists in the cache, then return it immediately
	if client.cache != nil {
		if cachedValue := client.cache.Get(uri); cachedValue != nil {
			return NewMap(cachedValue, client), nil
		}
	}

	// Try to load-and-parse the value from the remote server
	result := make(map[string]any)

	transaction := remote.Get(uri).
		Accept(client.mimeType).
		Response(&result, nil)

	if err := transaction.Send(); err != nil {
		return NewZero(), derp.Wrap(err, "jsonld.Client.Load", "Error loading JSON-LD document", uri)
	}

	// If we got a result, then cache it for later
	if client.cache != nil {
		client.cache.Set(uri, result)
	}

	// Return in triumph
	return NewMap(result, client), nil
}
