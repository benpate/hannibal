package streams

type Option func(*Document)

// WithClient option sets the HTTP client that can load remote documents if necessary
func WithClient(client Client) Option {
	return func(doc *Document) {
		if client == nil {
			doc.client = NewDefaultClient()
		} else {
			doc.client = client
		}
	}
}

// WithMeta option sets metadata in the document
func WithMeta(name string, value any) Option {
	return func(doc *Document) {
		doc.metadata[name] = value
	}
}
