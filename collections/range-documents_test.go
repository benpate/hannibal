//go:build localonly

package collections

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/davecgh/go-spew/spew"
)

func TestDocuments(t *testing.T) {

	doc := streams.NewDocument("https://mastodon.social/@benpate")
	outbox := doc.Outbox()

	items := RangeDocuments(outbox)

	index := 1
	for item := range items {
		spew.Dump(index)
		spew.Dump(item.Published())
		index++

		if index > 100 {
			break // okay, we get it.. you can load lots of documents...
		}
	}
}
