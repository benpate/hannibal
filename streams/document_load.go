package streams

import (
	"github.com/benpate/derp"
	"github.com/benpate/remote"
)

// AsObject retrieves a JSON-LD document from a remote server, parses is, and returns a Document object.
func (document Document) AsObject() (Document, error) {

	const location = "hannibal.streams.Document.AsObject"

	uri := document.ID()

	switch document.value.(type) {

	case map[string]any:
		return document, nil

	case []any:
		return document.Head().AsObject()

	case string:

		// If the value exists in the cache, then return it immediately
		if document.cache != nil {
			if cachedValue := document.cache.Get(uri); cachedValue != nil {
				return NewDocument(cachedValue, document.cache), nil
			}
		}

		// Try to load-and-parse the value from the remote server
		result := make(map[string]any)

		transaction := remote.Get(uri).
			Accept("application/activity+json").
			Response(&result, nil)

		if err := transaction.Send(); err != nil {
			return NilDocument(), derp.Wrap(err, location, "Error loading JSON-LD document", uri)
		}

		// If we got a result, then cache it for later
		if document.cache != nil {
			document.cache.Set(uri, result)
		}

		// Return in triumph
		return NewDocument(result, document.cache), nil
	}

	return NilDocument(), derp.NewInternalError(location, "Document type is invalid", document.Value())
}
