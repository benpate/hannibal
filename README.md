# Hannibal

<img src="https://github.com/benpate/hannibal/raw/main/meta/logo.jpg" style="width:100%; display:block; margin-bottom:20px;"  alt="Oil painting titled: Hannibal in the Alps, by R.B. Davis">

Hannibal is an experimental ActivityPub library for Go. It's goal is to be a robust, idiomatic, and thoroughly documented ActivityPub implementation fits into your application without any magic or drama.

There are other packages/frameworks out there that are more complete and mature. So please check out [go-fed](https://github.com/go-fed) and [go-ap](https://github.com/go-ap) before trying this.


## Packages
Like the ActivityPub spec itself, Hannibal is broken into several layers:

### pub - ActivityPub client/server
https://www.w3.org/TR/activitypub/

This is not an ActivityPub framework, but a simple library that easily plugs into your existing app.  Add ActivityPub behaviors to your existing handlers, and send ActivityPub messages to 

### vocab - ActivityStreams Vocabulary
https://www.w3.org/TR/activitystreams-vocabulary/

The `vocab` package includes the standard ActivityStream vocabulary, including names of actions, objects and properties used in ActivityPub. 

### streams - ActivityStreams data structures
https://www.w3.org/TR/activitystreams-core/

The `streams` package contains common data structures defined in the ActivityStreams spec, notably definitions for: `Document`, `Collection`, `OrderedCollection`, `CollectionPage`, and `OrderedCollectionPage`.  These are used by ActivityPub to send and receive multiple records in one HTTP request.

This package also includes a lightweight wrapper around generic data structures (like `map[string]any` and `[]any`) that makes it easy to access data structures within an ActivityStreams/JSON-LD document.

### sigs - HTTP Signatures and Digests
https://datatracker.ietf.org/doc/draft-ietf-httpbis-message-signatures

The `sigs` package creates and verifies HTTP signatures and Digests.