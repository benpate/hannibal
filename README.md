# Hannibal

<figure style="margin:0px;">
<img src="https://github.com/benpate/hannibal/raw/main/meta/logo.jpg">
<figcaption style="font-size:10px; text-align:right;"><i>Hannibal In The Alps</i> by <a href="https://en.wikipedia.org/wiki/Richard_Barrett_Davis">R.B. Davis</a>.</figcaption>
</figure>
<br>

Hannibal is an experimental ActivityPub library for Go. It's goal is to be a robust, idiomatic, and thoroughly documented ActivityPub implementation fits into your application without any magic or drama.

There are other packages/frameworks out there that are more complete and mature. So please check out [go-fed](https://github.com/go-fed) and [go-ap](https://github.com/go-ap) before trying this.


## Packages
Like the ActivityPub spec itself, Hannibal is broken into several layers:

### pub - ActivityPub client/server
https://www.w3.org/TR/activitypub/

This is not an ActivityPub framework, but a simple library that easily plugs into your existing app.  Add ActivityPub behaviors to your existing handlers, and send ActivityPub messages to 

### vocab - ActivityStreams Vocabulary
https://www.w3.org/TR/activitystreams-vocabulary/

This package includes the standard ActivityStream vocabulary, including names of actions, objects and properties used in ActivityPub. 

### streams - ActivityStreams data structures
https://www.w3.org/TR/activitystreams-core/

The stream package contains common data structures defined in the ActivityStreams spec, notably definitions for: `Collection`, `OrderedCollection`, `CollectionPage`, and `OrderedCollectionPage`.  These are used by ActivityPub to send and receive multiple records in one HTTP request.

### jsonld - JSON-LD reader/writer
https://json-ld.org

This package is a lightweight wrapper around generic data structures like `map[string]any` and `[]any` that simplifies the messy job of accessing the many data structures within JSON-LD.  If you're looking for a real, rigorous JSON-LD implementation, you should check out https://github.com/piprate/json-gold
