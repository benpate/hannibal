package streams

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
)

// Iterator allows transparent iteration over a Collection, OrderedCollection, and their corresponding pages (if present).
type Iterator struct {
	document Document
	counter  int
	err      error
}

// NewIterator creates a new Iterator object from a Document or Document URI.
func NewIterator(uri Document) (Iterator, error) {

	// If the uri is not an actual document, then load it now.
	document, err := uri.Load()

	if err != nil {
		return Iterator{}, derp.Wrap(err, "hannibal.streams.NewIterator", "Error retrieving document", document)
	}

	return Iterator{
		document: document,
		counter:  0,
		err:      nil,
	}, nil
}

// TotalItems returns the total number of items in the collection.
func (it *Iterator) TotalItems() int {
	return it.document.TotalItems()
}

// HasNext verifies that there is at least one more item remaining in the collection.
func (it *Iterator) HasNext() bool {

	// If we still have items available in the current page
	if items := it.getItems(); it.counter < len(items) {
		return true
	}

	// Otherwise, try to load the next page of results
	var page Document
	var err error

	switch it.document.Type() {

	// Collections point to the "first" page of items
	case vocab.CoreTypeCollection, vocab.CoreTypeOrderedCollection:
		page, err = it.document.First().Load()

	// CollectionPages point to the "next" page of items
	case vocab.CoreTypeCollectionPage, vocab.CoreTypeOrderedCollectionPage:
		page, err = it.document.Next().Load()

	}

	// Handle Errors
	if err != nil {
		it.err = derp.Wrap(err, "hannibal.streams.Iterator.HasNext", "Error loading next page of results from..", it.document.Value())
		return false
	}

	// If the document itself is empty, then ew, gross.. fail.
	if page.IsNil() {
		return false
	}

	// Otherwise, update the iterator
	it.document = page
	it.counter = 0

	// To be valid, the next page of results needs at least one item in it.
	items := it.getItems()
	return len(items) > 0
}

// Next returns the next Document in the Collection.
func (it *Iterator) Next() Document {

	items := it.getItems()

	// If we already have a document in memory with items in it, then it's easy to return.
	// When using .HasNext() then this should ALWAYS be the case.
	if it.counter < len(items) {
		result := items[it.counter]
		it.counter = it.counter + 1
		return it.document.sub(result)
	}

	// Call HasNext to see if we really DO have a next page or not.
	if it.HasNext() {
		return it.Next()
	}

	// Iterator has no more items in it, and you're a bad person for calling this.
	it.err = derp.NewInternalError("hannibal.streams.Iterator.GetNext", "Must call .HasNext() before calling .GetNext()", it.document.Value())
	return NilDocument()
}

// Error returns any error that occurred during the iteration process.
func (it *Iterator) Error() error {
	return it.err
}

// getItems retrieves the current page of items from the document
func (it *Iterator) getItems() []any {

	if items := it.document.Items(); !items.IsNil() {
		return items.Array()
	}

	return make([]any, 0)
}
