//gobuild:localonly

package collections

import (
	"sort"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/rosetta/channel"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestPoast(t *testing.T) {

	actor, err := streams.NewDocument("https://poa.st/users/benpate").Load()
	require.Nil(t, err)

	outbox, err := actor.Outbox().Load()
	require.Nil(t, err)

	/*
		pages := Pages(outbox, nil)

		for page := range pages {
			spew.Dump(page.Value())
		}
	*/

	done := make(chan struct{})
	documents := Documents(outbox, done)           // start reading documents from the outbox
	documents = channel.Limit(12, documents, done) // Limit to last 12 documents

	/*
		for document := range documents {
			spew.Dump(document.Value())
		}
	*/

	documentsSlice := channel.Slice(documents) // Convert the channel into a slice

	// Sort the collection chronologically so that they're imported in the correct order.
	sort.Slice(documentsSlice, func(a int, b int) bool {
		return documentsSlice[a].Published().Before(documentsSlice[b].Published())
	})

	spew.Dump(len(documentsSlice))

	for _, document := range documentsSlice {
		spew.Dump("=======", document.ID(), document.Published())
	}
}
