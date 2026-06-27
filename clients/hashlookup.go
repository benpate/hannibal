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

	// Otherwise, try to load the baseURL and find the hash inside that document.
	// Spread options... so they pass through individually; passing the slice
	// as one argument would nest it and silently drop caller options.
	result, err := client.innerClient.Load(baseURL, options...)

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
	return streams.NilDocument(), derp.NotFound("ashash.Client.Load", "Hash value not found in document", baseURL, hash, result.Value())
}

func (client HashLookup) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

// Delete removes a document from the underlying client's cache.
func (client HashLookup) Delete(documentID string) error {
	return client.innerClient.Delete(documentID)
}

// SetRootClient passes the top-level client down to the underlying client, so
// stacked clients that make recursive calls resolve through the whole chain.
func (client HashLookup) SetRootClient(rootClient streams.Client) {
	client.innerClient.SetRootClient(rootClient)
}

// Verify that HashLookup satisfies the streams.Client interface.
var _ streams.Client = HashLookup{}
