// Package collection is an HTTP interface for GET-ing activities from
// an ActivityPub collection.
package collection

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/ranges"
	"github.com/labstack/echo/v4"
)

// maxItemsPerPage is the maximum number of items that will be returned in a single page of an OrderedCollection
const maxItemsPerPage = 60

// Serve generates a data for the collection that can be retured via HTTP,
// returning ActivityStreams OrderedCollection or OrderedCollectionPage of outbound activities.
// The `storage` parameter is assumed to already contain any necessary filtering information
// such as Actor permissions and "after" cursors.
func Serve(ctx echo.Context, collectionID string, countFunc CounterFunc, iteratorFunc IteratorFunc, options ...Option) error {

	config := NewConfig(options...)

	// If we have an "after" parameter, then return the OrderedCollectionPage
	// corresponding to all activities after the provided ID
	if after := ctx.QueryParam("after"); after != "" {
		return serveOrderedCollectionPage(ctx, collectionID, iteratorFunc, after)
	}

	// Otherwise, return the OrderedCollection container
	return serveOrderedCollection(ctx, collectionID, countFunc, iteratorFunc, config)
}

func serveOrderedCollection(ctx echo.Context, collectionID string, countFunc CounterFunc, iteratorFunc IteratorFunc, config Config) error {

	const location = "collection.serveOrderedCollection"

	// Count the total number of activities in this outbox
	totalItems, err := countFunc()

	if err != nil {
		return derp.Wrap(err, location, "Unable to count activities in outbox storage")
	}

	result := mapof.Any{
		vocab.AtContext:          vocab.ContextTypeActivityStreams,
		vocab.PropertyType:       vocab.CoreTypeOrderedCollection,
		vocab.PropertyID:         collectionID,
		vocab.PropertyTotalItems: totalItems,
	}

	if config.SSEEndpoint != "" {
		result[vocab.PropertyEventStream] = config.SSEEndpoint
	}

	if totalItems == 0 {
		return serveJSON(ctx, http.StatusOK, result)
	}

	result[vocab.PropertyFirst] = collectionID + "?after=FIRST"
	return serveJSON(ctx, http.StatusOK, result)
}

// serveOrderedCollectionPage serves a single page of activities from the OrderedCollection
// starting after the provided activity ID.
func serveOrderedCollectionPage(ctx echo.Context, collectionID string, iteratorFunc IteratorFunc, after string) error {

	const location = "collection.serveOrderedCollectionPage"

	// Remove the magic "FIRST" value (if present) to send to the getItems() function
	afterParam := iif(after == "FIRST", "", after)

	// Get the activities in the collection
	orderedItems, lastID, err := getItems(iteratorFunc, afterParam)

	if err != nil {
		return derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	// Build the OrderedCollectionPage response
	result := mapof.Any{
		vocab.AtContext:            vocab.ContextTypeActivityStreams,
		vocab.PropertyType:         vocab.CoreTypeOrderedCollectionPage,
		vocab.PropertyID:           collectionID + "?after=" + after,
		vocab.PropertyPartOf:       collectionID,
		vocab.PropertyOrderedItems: orderedItems,
	}

	// Set the "next" after URL (if applicable)
	if len(orderedItems) >= maxItemsPerPage {
		result[vocab.PropertyNext] = collectionID + "?after=" + lastID
	}

	// Done!
	return ctx.JSON(http.StatusOK, result)
}

// getItems retrieves all items from the collection after the provided "lastID".
func getItems(iteratorFunc IteratorFunc, after string) (items []any, lastID string, err error) {

	const location = "collection.getItems"

	// Initialize the result so that we can return an empty array if there are no items in the collection
	items = make([]any, 0)

	// Get an iterator over all (matching) activities in the Storage adapter
	activities, err := iteratorFunc(after)

	if err != nil {
		return nil, "", derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	// Do not retrieve more than 60 items
	activities = ranges.Limit(maxItemsPerPage, activities)

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

func serveJSON(ctx echo.Context, statusCode int, data any) error {

	if hannibal.IsActivityPubContentType(ctx.Request().Header.Get("Accept")) {
		return ctx.JSON(statusCode, data)
	}

	return ctx.JSONPretty(statusCode, data, "    ")
}
