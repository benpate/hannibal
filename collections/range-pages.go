// Package collections provides tools for traversing ActivityStreams collections,
// represented as streams.Documents, including pagination.
package collections

import (
	"iter"

	"github.com/benpate/hannibal/streams"
)

// RangePages returns an iterator over every page of a paged collection.
func RangePages(collection streams.Document) iter.Seq[streams.Document] {

	return func(yield func(streams.Document) bool) {

		// emptyPage is used to prevent WriteFreely-style infinite loops
		var emptyPage bool

		// If this is a collection header, then move to the first page of results.
		// LoadLink (not Load) is required because "first"/"next" may be either a URL
		// string (which must be fetched) OR an inline/embedded page with no "id" (which
		// must be used as-is). Calling Load() on an inline page returns NilDocument
		// because it has no ID to fetch, which would silently drop the whole collection.
		if firstPage := collection.First(); firstPage.NotNil() {
			collection = firstPage.LoadLink()
		}

		// As long as we have a valid collection...
		for collection.NotNil() {

			// Send the collection to the caller
			if !yield(collection) {
				return
			}

			// Move to the next page in the collection (if available), loading it from
			// a URL when necessary. See the note above on why this is LoadLink, not Load.
			collection = collection.Next().LoadLink()

			// If this document is an empty page, then try to prevent
			// WriteFreely-style infinite loops.
			if collection.Items().Len() == 0 {

				// If we've already seen ONE empty page, then exit.
				if emptyPage {
					return
				}

				// Otherwise, set the emptyPage flag so we don't loop indefinitely.
				emptyPage = true
			}
		}
	}
}
