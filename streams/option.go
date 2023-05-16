package streams

type Option func(*Document)

// WithClient option sets a
func WithClient(client Client) Option {
	return func(doc *Document) {
		if client == nil {
			doc.client = NewDefaultClient()
		} else {
			doc.client = client
		}
	}
}
