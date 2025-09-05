package collections

import (
	"iter"

	"github.com/benpate/hannibal/streams"
)

func RangeDocuments(collection streams.Document) iter.Seq[streams.Document] {

	return func(yield func(streams.Document) bool) {

		// Loop through every page in the collection
		for page := range RangePages(collection) {

			// Loop through all items in the page
			for items := page.Items(); items.NotNil(); items = items.Tail() {

				if !yield(items.Head()) {
					return
				}
			}
		}
	}
}
