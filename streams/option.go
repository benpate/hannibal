package streams

import "net/http"

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

// WithHeader option sets the HTTP header that was returned by a remote HTTP request.
func WithHeader(header http.Header) Option {
	return func(doc *Document) {
		doc.header = header
	}
}

// WithLocales option sets one or more locales to use when working with language maps
func WithLocales(locales ...string) Option {
	return func(doc *Document) {
		doc.locales = locales
	}
}
