package streams

// Client represents an HTTP client (or facades in front of one)
// that can load a JSON-LD document from a remote server. A Client
// is injected into each streams.Document record so that the
// Document can load additional linked data as needed.
type Client interface {

	// Load returns a Document representing the specified URI.
	Load(uri string, options ...any) (Document, error)
}
