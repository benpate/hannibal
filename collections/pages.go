// Package iterator provides utilities for iterating through remote collections (represented as streams.Documents)
package collections

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/rosetta/channel"
)

func Pages(collection streams.Document, done <-chan struct{}) <-chan streams.Document {

	var err error
	result := make(chan streams.Document, 1)

	go func() {

		// emptyPage is used to prevent WriteFreely-style infinite loops
		var emptyPage bool

		defer close(result)

		// If this is a collection header, then try to load the first page of results
		if firstPage := collection.First(); firstPage.NotNil() {
			collection, err = firstPage.Load()

			if err != nil {
				derp.Report(derp.Wrap(err, "hannibal.collections.Iterator", "Error loading first page", collection))
				return
			}
		}

		// As long as we have a valid collection...
		for collection.NotNil() {

			// Breakpoint for cancellation
			if channel.Closed(done) {
				return
			}

			// Send the collection to the caller
			result <- collection

			// Look for the next page in the collection (if available)
			collection = collection.Next()

			// Try to load it and continue the loop.
			collection, err = collection.Load()

			if err != nil {
				derp.Report(derp.Wrap(err, "hannibal.collections.Iterator", "Error loading first page", collection))
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
	}()

	return result
}
