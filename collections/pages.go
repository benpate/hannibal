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

		defer close(result)

		// If this is a collection header, then try to load the first page of results
		if firstPage := collection.First(); firstPage.NotNil() {
			collection, err = firstPage.Load()

			if err != nil {
				// nolint:errcheck
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
				// nolint:errcheck
				derp.Report(derp.Wrap(err, "hannibal.collections.Iterator", "Error loading first page", collection))
				return
			}
		}
	}()

	return result
}

func PagesReverse(collection streams.Document, done <-chan struct{}) <-chan streams.Document {

	var err error
	result := make(chan streams.Document, 1)

	go func() {

		defer close(result)

		// If this is a collection header, then try to load the first page of results
		if lastPage := collection.Last(); lastPage.NotNil() {
			collection, err = lastPage.Load()

			if err != nil {
				// nolint:errcheck
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
			collection = collection.Prev()

			// Try to load it and continue the loop.
			collection, err = collection.Load()

			if err != nil {
				// nolint:errcheck
				derp.Report(derp.Wrap(err, "hannibal.collections.Iterator", "Error loading first page", collection))
				return
			}
		}
	}()

	return result
}
