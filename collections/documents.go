package collections

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/rosetta/channel"
)

func Documents(collection streams.Document, done <-chan struct{}) <-chan streams.Document {

	pages := Pages(collection, done)
	result := make(chan streams.Document, 1)

	go func() {

		defer close(result)

		for page := range pages {

			// Loop through all items in the page
			for items := page.Items(); items.NotNil(); items = items.Tail() {

				// Breakpoint for cancellation
				if channel.Closed(done) {
					return
				}

				// Return the next item and move forward one step.
				result <- items.Head()
			}
		}
	}()

	return result
}

func DocumentsReverse(collection streams.Document, done <-chan struct{}) <-chan streams.Document {
	pages := PagesReverse(collection, done)

	result := make(chan streams.Document, 1)

	go func() {

		defer close(result)

		for page := range pages {

			// Retrieve all items in the collection
			items := page.Items()

			for index := items.Len() - 1; index >= 0; index-- {

				// Breakpoint for cancellation
				if channel.Closed(done) {
					return
				}

				// Return the next item and move (backward) one step
				result <- items.At(index)
			}
		}
	}()

	return result
}
