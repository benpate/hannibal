# Hannibal / collections

This package provides channel-based tools for traversing [ActivityStreams](https://www.w3.org/TR/activitystreams-core/) collections.  

```go
// Retrieve a collection from the interwebs
outboxCollection := streams.NewDocument("https://your-website/@your-actor/outbox")

// Create a channel to iterate over the collection
documentChannel := collections.Documents(outboxCollection, context.TODO().Done())

// Yep. That's all there is. Now get to work.
for document := range documentChannel {
	// do stuff.
}
```

### Traversing All Documents in a Collection

The `Documents()` function is probably the only function you'll need from this package.  It returns a 
channel of all documents in a collection. It works with both 
[`Collections`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-collection), 
and [`OrderedCollections`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollection),
regardless of whether all documents are included directly in the collection, or if they are spread
across multiple pages.

### Traversing All Pages in a Collection

The `Pages` function returns a channel of all pages in a collection.  It can traverse both 
[`CollectionPage`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-collectionpage)s 
and [`OrderedCollectionPage`](https://www.w3.org/TR/activitystreams-vocabulary/#dfn-orderedcollectionpage)s


## Additional Tools

The [rosetta channel](https://github.com/benpate/rosetta/tree/main/channel) package includes a number of functions for manipulating channels.  For instance, if you don't want to read an actor's entire outbox, you might limit the results like this:

```go
// Retrieve a collection from the interwebs
outboxCollection := streams.NewDocument("https://your-website/@your-actor/outbox")

// Create a "done" channel to cancel iterating over the collection
var done chan struct{}

// Create a channel to iterate over the collection
documentChannel := collections.Documents(outboxCollection, done)

// Limit will use the "done" channel to cancel iteration once we reach the limit
limitedChannel := channel.Limit(10, documentChannel, done)

// Now just iterate.  The limitedChannel will close after 10 documents.
for document := range limtedChannel {
	// do stuff.
}
```
