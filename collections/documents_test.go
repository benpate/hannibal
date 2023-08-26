//go:build localonly

package collections

import (
	"context"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/davecgh/go-spew/spew"
)

func TestDocuments(t *testing.T) {

	doc := streams.NewDocument("https://mastodon.social/@benpate")
	outbox := doc.Outbox()

	items := Documents(outbox, context.TODO().Done())

	index := 1
	for item := range items {
		spew.Dump(index)
		spew.Dump(item.Published())
		index++
	}
}
