package collections

import "github.com/benpate/hannibal/streams"

// NewIterator is API-sugar for collections.Documents() iterator.
func NewIterator(collection streams.Document, options ...IteratorOption) <-chan streams.Document {

	done := make(chan struct{})
	result := Documents(collection, done)

	for _, option := range options {
		result = option(result, done)
	}

	return result
}
