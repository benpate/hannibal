package collections

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// CountItems returns the number of items in the given collection. If the
// collection includes a `TotalItems` value, then this is returned. Otherwise,
// this function will iterate through all pages in the collection to count the items.
func CountItems(collection streams.Document) (int, error) {

	const location = "hannibal.collections.CountItems"

	// If the collection does not exist, then there are no items.
	if collection.IsNil() {
		return 0, nil
	}

	defer func() {
		if err := recover(); err != nil {
			derp.Report(derp.Internal(location, "Recovered error", err, collection.ID()))
		}
	}()

	// If this is not already a "map" then load the complete document.
	collection = collection.LoadLink()

	// If the collection already reports the `TotalItems` count, then just use that.
	// This is the best case scenario.
	if totalItems := collection.TotalItems(); totalItems > 0 {
		return totalItems, nil
	}

	// Otherwise, we have to do your work for you.  Let's start paging and counting...
	var result int

	// Retrieve each page in the collection and count the number of items it contains.
	for page := range RangePages(collection) {

		// Get the number of items in this page
		count := page.Items().Len()

		// If this page is empty, then we don't need to load any more pages
		if count == 0 {
			break
		}

		// Increment the total count
		result += count
	}

	// Return the number of items found
	return result, nil
}
