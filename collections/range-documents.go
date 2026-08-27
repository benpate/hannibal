package collections

import (
	"iter"

	"github.com/benpate/hannibal/streams"
)

// RangeDocuments returns an iterator over every Document in the collection, following pagination.
func RangeDocuments(collection streams.Document, options ...Option) iter.Seq[streams.Document] {

	config := newRangeConfig(options...)

	return func(yield func(streams.Document) bool) {

		documentCount := 0

		// Loop through every page in the collection (forwarding traversal options)
		for page := range RangePages(collection, options...) {

			// Loop through all items in the page
			for items := page.Items(); items.NotNil(); items = items.Tail() {

				// RULE: Cap the total number of documents.  The page cap alone still
				// admits a remote server stuffing each page with an enormous "items"
				// array, so this bounds the total work handed to the caller.
				if documentCount >= config.maxDocuments {
					return
				}

				if !yield(items.Head()) {
					return
				}

				documentCount++
			}
		}
	}
}
