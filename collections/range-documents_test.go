//go:build localonly

package collections

import (
	"testing"

	"github.com/benpate/hannibal/streams"
)

func TestDocuments(t *testing.T) {

	doc := streams.NewDocument("https://mastodon.social/@benpate")
	outbox := doc.Outbox()

	items := RangeDocuments(outbox)

	index := 1
	for item := range items {
		t.Log(item.Published())
		index++

		if index > 100 {
			break // okay, we get it.. you can load lots of documents...
		}
	}
}
