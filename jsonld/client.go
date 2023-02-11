package jsonld

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/cache"
	"github.com/benpate/remote"
)

// Client wraps http transactions to load remote JSON-LD documents.
type Client struct {
	cache Cache
}

// New creates a new Client object, which can be used to load remote JSON-LD documents.
func New(cache Cache) Client {
	return Client{
		cache: cache,
	}
}

func NewDefaultClient() Client {
	return Client{
		cache: cache.NewDefaultCache(),
	}
}

// NewReader returns a Reader object for the specified value.
func (client *Client) NewReader(value any) Reader {

	switch typed := value.(type) {

	case map[string]any:
		return NewMap(typed, client)

	case []any:
		return NewSlice(typed, client)

	case string:
		return NewString(typed, client)

	case bool:
		return NewBool(typed)

	case int:
		return NewInt(typed)

	case int64:
		return NewInt(int(typed))

	case float64:
		return NewFloat(typed)

	}

	return NewZero()
}

// Load retrieves a JSON-LD document from a remote server, parses is, and returns a Reader object.
func (client *Client) Load(uri string) (Reader, error) {

	// If the value exists in the cache, then return it immediately
	if client.cache != nil {
		if cachedValue := client.cache.Get(uri); cachedValue != nil {
			return NewMap(cachedValue, client), nil
		}
	}

	// Try to load-and-parse the value from the remote server
	result := make(map[string]any)

	transaction := remote.Get(uri).
		Accept("application/activity+json").
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
