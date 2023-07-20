//go:build localonly

package collections

import (
	"context"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/davecgh/go-spew/spew"
)

func TestPages(t *testing.T) {

	doc := streams.NewDocument("https://mastodon.social/@benpate")
	outbox := doc.Outbox()

	pages := Pages(outbox, context.TODO().Done())

	index := 1
	for page := range pages {
		spew.Dump(index)
		spew.Dump(page.ID())
		index++
	}
}

func TestPagesReverse(t *testing.T) {

	doc := streams.NewDocument("https://mastodon.social/@benpate")
	outbox := doc.Outbox()

	pages := PagesReverse(outbox, context.TODO().Done())

	index := 1
	for page := range pages {
		spew.Dump(index)
		spew.Dump(page.ID())
		index++
	}
}
