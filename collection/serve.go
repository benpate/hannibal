// Package collection is an HTTP interface for GET-ing activities from
// an ActivityPub collection.
package collection

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// Serve generates a data for the collection that can be retured via HTTP,
// returning ActivityStreams OrderedCollection or OrderedCollectionPage of outbound activities.
// The `storage` parameter is assumed to already contain any necessary filtering information
// such as Actor permissions and "after" cursors.
func Serve(storage Storage, after string) (mapof.Any, error) {

	// If we have an "after" parameter, then return the OrderedCollectionPage
	// corresponding to all activities after the provided ID
	if after != "" {
		return serveOrderedCollectionPage(storage, after)
	}

	// Otherwise, return the OrderedCollection container
	return serveOrderedCollection(storage)
}

func serveOrderedCollection(storage Storage) (mapof.Any, error) {

	const location = "collection.serveOrderedCollection"

	// Count the total number of activities in this outbox
	totalItems, err := storage.TotalItems()

	if err != nil {
		return mapof.NewAny(), derp.Wrap(err, location, "Unable to count activities in outbox storage")
	}

	// If there are more than 60 items in the collection, then
	// don't include them all here.  Instead, provide a "first" link
	// to the first page of results.
	if totalItems > 60 {

		return mapof.Any{
			vocab.AtContext:          vocab.ContextTypeActivityStreams,
			vocab.PropertyType:       vocab.CoreTypeOrderedCollection,
			vocab.PropertyID:         storage.ID(),
			vocab.PropertyTotalItems: totalItems,
			vocab.PropertyFirst:      storage.ID() + "?after=0",
		}, nil
	}

	// Retrieve all items from the storage adapter
	orderedItems, _, err := getItems(storage, "0")

	if err != nil {
		return mapof.NewAny(), derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	return mapof.Any{
		vocab.AtContext:            vocab.ContextTypeActivityStreams,
		vocab.PropertyType:         vocab.CoreTypeOrderedCollection,
		vocab.PropertyID:           storage.ID(),
		vocab.PropertyTotalItems:   totalItems,
		vocab.PropertyOrderedItems: orderedItems,
	}, nil
}

// serveOrderedCollectionPage serves a single page of activities from the OrderedCollection
// starting after the provided activity ID.
func serveOrderedCollectionPage(storage Storage, after string) (mapof.Any, error) {

	const location = "collection.serveOrderedCollectionPage"

	// Get the activities in the collection
	orderedItems, lastID, err := getItems(storage, after)

	if err != nil {
		return mapof.NewAny(), derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	// Build the OrderedCollectionPage response
	result := mapof.Any{
		vocab.AtContext:            vocab.ContextTypeActivityStreams,
		vocab.PropertyType:         vocab.CoreTypeOrderedCollectionPage,
		vocab.PropertyID:           storage.ID() + "?after=" + after,
		vocab.PropertyPartOf:       storage.ID(),
		vocab.PropertyFirst:        storage.ID() + "?after=0",
		vocab.PropertyOrderedItems: orderedItems,
	}

	// Set the "next" after URL (if applicable)
	if lastID != "" {
		result[vocab.PropertyNext] = storage.ID() + "?after=" + lastID
	}

	// Done!
	return result, nil
}

// getItems retrieves all items from the collection after the provided "lastID".
func getItems(storage Storage, after string) (items []any, lastID string, err error) {

	const location = "collection.getItems"

	// Get an iterator over all (matching) activities in the Storage adapter
	activities, err := storage.Iterator(after)

	if err != nil {
		return nil, "", derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	for activity := range activities {
		lastID = activity.GetString(vocab.PropertyID)
		items = append(items, activityValue(activity))
	}

	return items, lastID, nil
}

// activityValue returns either the activity ID as a string if only this value is present,
// or a full mapof.Any if more than just the ID is present
func activityValue(activity mapof.Any) any {

	if activityID := activity.GetString(vocab.PropertyID); (len(activity) == 1) && (activityID != "") {
		return activityID
	}

	return activity
}
