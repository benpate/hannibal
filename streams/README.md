## Hannibal / Streams

This package implements the [ActivityStreams 2.0](https://www.w3.org/TR/activitystreams-core/) 
specifications in Hannibal.  Specifically, it provides an overly simplistic view of JSON-LD 
documents, contexts, and the various types of ActivityStreams collections.

This is not a rigorous implementation of JSON-LD.  Instead, it is an easy way to navigate 
well-formatted JSON-LD documents and to iterate over their contents, as well as loading additional documents from the web when necessary.

```go
// Load a document directly from its URL
document, err := streams.Load("https://your.activitypub.server/@documentId")

document.ID() // returns a string value
document.Content() // returns the content property
document.Published() // returns a time value

// AttributedTo could be many things.. a single value, a link, or
// an array of values and links. Let's make all of that easier.
authors := document.AttributedTo() 

// You could just read the first author directly
authors.Name() // returns the 'Name' string
authors.ID() // returns the 'ID' string

// Or you can use it as an iterator
for authors := document.AttributedTo ; !authors.IsNil() ; authors = authors.Tail() {
	authors.Value() // returns the whole value from the array
}
```