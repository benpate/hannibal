package streams

type Client interface {
	// LoadActor returns a Document representing the Actor at the specified URI, which must
	// contain an "outbox" property.
	LoadActor(uri string) (Document, error)

	// LoadDocument returns a Document representing the specified URI
	LoadDocument(uri string, defaultValue map[string]any) (Document, error)
}
