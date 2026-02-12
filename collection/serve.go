// Package collection is an HTTP interface for GET-ing activities from
// an ActivityPub collection.
package collection

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/labstack/echo/v4"
)

// Serve generates a data for the collection that can be retured via HTTP,
// returning ActivityStreams OrderedCollection or OrderedCollectionPage of outbound activities.
// The `storage` parameter is assumed to already contain any necessary filtering information
// such as Actor permissions and "after" cursors.
func Serve(ctx echo.Context, idFunc IdentifierFunc, countFunc CounterFunc, iteratorFunc IteratorFunc) error {

	// If we have an "after" parameter, then return the OrderedCollectionPage
	// corresponding to all activities after the provided ID
	if after := ctx.QueryParam("after"); after != "" {
		return serveOrderedCollectionPage(ctx, idFunc, countFunc, iteratorFunc, after)
	}

	// Otherwise, return the OrderedCollection container
	return serveOrderedCollection(ctx, idFunc, countFunc, iteratorFunc)
}

func serveOrderedCollection(ctx echo.Context, idFunc IdentifierFunc, countFunc CounterFunc, iteratorFunc IteratorFunc) error {

	const location = "collection.serveOrderedCollection"

	// Count the total number of activities in this outbox
	totalItems, err := countFunc()

	if err != nil {
		return derp.Wrap(err, location, "Unable to count activities in outbox storage")
	}

	result := mapof.Any{
		vocab.AtContext:          vocab.ContextTypeActivityStreams,
		vocab.PropertyType:       vocab.CoreTypeOrderedCollection,
		vocab.PropertyID:         idFunc(),
		vocab.PropertyTotalItems: totalItems,
	}

	if totalItems == 0 {
		return serveJSON(ctx, http.StatusOK, result)
	}

	// If there are more than 60 items in the collection, then
	// don't include them all here.  Instead, provide a "first" link
	// to the first page of results.
	if totalItems > 60 {
		result[vocab.PropertyFirst] = idFunc() + "?after=FIRST"
		return serveJSON(ctx, http.StatusOK, result)
	}

	// Retrieve all items from the storage adapter
	orderedItems, _, err := getItems(iteratorFunc, "FIRST")

	if err != nil {
		return derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	result[vocab.PropertyOrderedItems] = orderedItems
	return serveJSON(ctx, http.StatusOK, result)
}

// serveOrderedCollectionPage serves a single page of activities from the OrderedCollection
// starting after the provided activity ID.
func serveOrderedCollectionPage(ctx echo.Context, idFunc IdentifierFunc, countFunc CounterFunc, iteratorFunc IteratorFunc, after string) error {

	const location = "collection.serveOrderedCollectionPage"

	// Get the activities in the collection
	orderedItems, lastID, err := getItems(iteratorFunc, after)

	if err != nil {
		return derp.Wrap(err, location, "Unable to retrieve outbox activities")
	}

	collectionID := idFunc()

	// Build the OrderedCollectionPage response
	result := mapof.Any{
		vocab.AtContext:            vocab.ContextTypeActivityStreams,
		vocab.PropertyType:         vocab.CoreTypeOrderedCollectionPage,
		vocab.PropertyID:           collectionID + "?after=" + after,
		vocab.PropertyPartOf:       collectionID,
		vocab.PropertyFirst:        collectionID + "?after=FIRST",
		vocab.PropertyOrderedItems: orderedItems,
	}

	// Set the "next" after URL (if applicable)
	if lastID != "" {
		result[vocab.PropertyNext] = collectionID + "?after=" + lastID
	}

	// Done!
	return ctx.JSON(http.StatusOK, result)
}

// getItems retrieves all items from the collection after the provided "lastID".
func getItems(iteratorFunc IteratorFunc, after string) (items []any, lastID string, err error) {

	const location = "collection.getItems"

	// Get an iterator over all (matching) activities in the Storage adapter
	activities, err := iteratorFunc(after)

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

func serveJSON(ctx echo.Context, statusCode int, data any) error {

	if hannibal.IsActivityPubContentType(ctx.Request().Header.Get("Accept")) {
		return ctx.JSON(statusCode, data)
	}

	return ctx.JSONPretty(statusCode, data, "    ")
}
