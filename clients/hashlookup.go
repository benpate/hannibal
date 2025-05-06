package clients

import (
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// HashLookup is a streams.Client wrapper that searches for hash values in a document.
type HashLookup struct {
	innerClient streams.Client
}

// NewHashLookup creates a fully initialized Client object
func NewHashLookup(innerClient streams.Client) HashLookup {
	return HashLookup{
		innerClient: innerClient,
	}
}

// Load retrieves a document from the underlying innerClient, then searches for hash values
// inside it (if required)
func (client HashLookup) Load(url string, options ...any) (streams.Document, error) {

	// Try to find a hash in the URL
	baseURL, hash, found := strings.Cut(url, "#")

	// If there is no hash, then proceed as is.
	if !found {
		return client.innerClient.Load(url, options...)
	}

	// Otherwise, try to load the baseURL and find the hash inside that document
	result, err := client.innerClient.Load(baseURL, options)

	if err != nil {
		return result, err
	}

	// Search all properties at the top level of the document (not recursive)
	// and scan through arrays (if present) looking for an ID that matches the original URL (base + hash)
	for _, key := range result.MapKeys() {
		for property := result.Get(key); property.NotNil(); property = property.Tail() {
			if property.ID() == url {
				return property, nil
			}
		}
	}

	// Not found.
	return streams.NilDocument(), derp.NotFoundError("ashash.Client.Load", "Hash value not found in document", baseURL, hash, result.Value())
}
