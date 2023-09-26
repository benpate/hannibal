package streams

import "net/http"

// HTTPHeader returns the http.Header object associated with this document
func (document Document) HTTPHeader() http.Header {
	return document.httpHeader
}

// SetHTTPHeader sets the http.Header object associated with this document
func (document *Document) SetHTTPHeader(httpHeader http.Header) {
	document.httpHeader = httpHeader
}
