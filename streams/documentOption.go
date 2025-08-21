package streams

import "net/http"

type DocumentOption func(*Document)

// WithClient option sets the HTTP client that can load remote documents if necessary
func WithClient(client Client) DocumentOption {
	return func(doc *Document) {
		if client == nil {
			doc.client = NewDefaultClient()
		} else {
			doc.client = client
		}
	}
}

// WithHTTPHeader attaches an HTTP header to the document
func WithHTTPHeader(httpHeader http.Header) DocumentOption {
	return func(doc *Document) {
		doc.httpHeader = httpHeader
	}
}

// WithMetadata attaches metadata to the document
func WithMetadata(metadata Metadata) DocumentOption {
	return func(doc *Document) {
		doc.Metadata = metadata
	}
}
