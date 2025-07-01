package streams

import (
	"iter"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
)

func (document Document) Range() iter.Seq[Document] {

	return func(yield func(Document) bool) {
		for ; document.NotNil(); document = document.Tail() {
			if !yield(document.Head()) {
				return
			}
		}
	}
}

// RangeIDs returns an iterator that yields the IDs of each element of a Document.
// If the Document is empty, it yields no values.
// If the Document is a single object, then it yields the ID of that object.
// If the Document is a list, then it yields the IDs of each object in the list
func (document Document) RangeIDs() iter.Seq[string] {

	return func(yield func(string) bool) {
		for ; document.NotNil(); document = document.Tail() {
			if !yield(document.Head().ID()) {
				return
			}
		}
	}
}

func (document Document) RangeMentions() iter.Seq[string] {
	return func(yield func(string) bool) {
		for tag := document.Tag(); tag.NotNil(); tag = tag.Tail() {

			if tag.Type() == vocab.LinkTypeMention {
				if !yield(tag.Href()) {
					return
				}
			}
		}
	}
}

func (document Document) RangeAddressees() iter.Seq[string] {

	return joinIterators(
		document.Actor().RangeIDs(),
		document.To().RangeIDs(),
		document.CC().RangeIDs(),
		document.BTo().RangeIDs(),
		document.BCC().RangeIDs(),
		document.RangeMentions(),

		// TODO: FEP-1b12: Group Federation: https://w3id.org/fep/1b12
		// TODO: FEP-7888: Demystifying the context property: https://w3id.org/fep/7888
		// TODO: FEP-7458: Using the replies collection: https://w3id.org/fep/7458
		// TODO: FEP-171b: Conversation Containers: http://w3id.org/fep/171b
		// TODO: FEP-f228: Backfilling conversations: https://w3id.org/fep/f228
	)
}

func (document Document) RangeInReplyTo() iter.Seq[string] {

	return func(yield func(string) bool) {

		inReplyTo := document.InReplyTo()

		if inReplyTo.IsNil() {
			return // Nothing to yield
		}

		inReplyToDocument, err := inReplyTo.Load()

		if err != nil {
			derp.Report(derp.Wrap(err, "streams.Document.RangeInReplyTo", "Error loading InReplyTo document", inReplyTo.ID()))
			return // Nothing to yield
		}

		for address := range inReplyToDocument.RangeAddressees() {
			if !yield(address) {
				return // Stop yielding
			}
		}
	}
}
