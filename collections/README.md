# Hannibal / collections

This package provides tools for traversing [ActivityStreams](https://www.w3.org/TR/activitystreams-core/) collections. It works with both [`Collections`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-collection) and [`OrderedCollections`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollection), whether the items are embedded directly in the collection or spread across multiple pages — following pagination automatically as it goes.

The traversal functions return Go 1.23 iterators (`iter.Seq[streams.Document]`), so you consume them with a plain `range` loop.

```go
// Retrieve a collection from the interwebs
outbox := streams.NewDocument("https://your-website/@your-actor/outbox")

// RangeDocuments iterates over every document in the collection,
// loading additional pages from the web as needed.
for document := range collections.RangeDocuments(outbox) {
	// do stuff
}
```

## Functions

- **`RangeDocuments(collection)`** — an iterator over every document in the collection, following pagination. This is the one you'll usually want.
- **`RangePages(collection)`** — an iterator over each *page* of a paged collection, rather than the individual documents.
- **`RangeDocumentsBefore(iterator, limit)`** — wraps a document iterator to stop after `limit` items published before the cursor. Use it to read just the start of a large collection.
- **`CountItems(collection)`** — returns the total number of items, walking every page if the collection does not report a `totalItems` count directly.

## Limiting Traversal

To read only part of a large collection, wrap the iterator with `RangeDocumentsBefore`, or break out of the `range` loop early — iteration stops and no further pages are fetched.

```go
outbox := streams.NewDocument("https://your-website/@your-actor/outbox")

// Read at most the 10 most recent documents
for document := range collections.RangeDocumentsBefore(collections.RangeDocuments(outbox), 10) {
	// do stuff
}
```
