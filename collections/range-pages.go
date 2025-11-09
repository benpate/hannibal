// Package iterator provides utilities for iterating through remote collections (represented as streams.Documents)
package collections

import (
	"iter"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

func RangePages(collection streams.Document) iter.Seq[streams.Document] {

	const location = "hannibal.collections.Pages"

	return func(yield func(streams.Document) bool) {

		var err error

		// emptyPage is used to prevent WriteFreely-style infinite loops
		var emptyPage bool

		// If this is a collection header, then try to load the first page of results
		if firstPage := collection.First(); firstPage.NotNil() {
			collection, err = firstPage.Load()

			if err != nil {
				derp.Report(derp.Wrap(err, location, "Unable to load first page", collection))
				return
			}
		}

		// As long as we have a valid collection...
		for collection.NotNil() {

			// Send the collection to the caller
			if !yield(collection) {
				return
			}

			// Look for the next page in the collection (if available)
			collection = collection.Next()

			// Try to load it and continue the loop.
			collection, err = collection.Load()

			if err != nil {
				derp.Report(derp.Wrap(err, location, "Unable to load first page", collection))
				return
			}

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
