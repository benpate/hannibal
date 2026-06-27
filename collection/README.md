# Hannibal / collection

Collection serves an ActivityStreams [`OrderedCollection`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollection)
(or a single [`OrderedCollectionPage`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollectionpage))
over HTTP. You supply two functions — one that counts the items and one that iterates them — and
`Serve` handles paging, content negotiation, and JSON-LD rendering.

```go
func handleOutbox(ctx echo.Context) error {

	// Returns the total number of items in the collection
	countFunc := func() (int64, error) {
		return myService.CountActivities(actorID)
	}

	// Returns an iterator over items, starting after the given cursor
	iteratorFunc := func(startIndex string) (iter.Seq[mapof.Any], error) {
		return myService.ActivitiesAfter(actorID, startIndex)
	}

	return collection.Serve(ctx, collectionID, countFunc, iteratorFunc,
		collection.WithAttributedTo(actorID),
	)
}
```

## Paging

`Serve` inspects the request's `after` query parameter:

- **no `after`** — returns the top-level `OrderedCollection` (with `first`/`last` page links and the total count).
- **`after=<id>`** — returns the `OrderedCollectionPage` of items following that cursor.

The `countFunc` and `iteratorFunc` are expected to already encode any filtering (actor permissions,
visibility, cursors); `Serve` does not add its own access control.

## Options

- `WithAttributedTo(id)` — sets the collection's `attributedTo` property.
- `WithAudience(id)` — sets the collection's `audience` property.
- `WithSSEEndpoint(url)` — advertises a Server-Sent Events endpoint for live updates.
