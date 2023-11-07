package collections

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/rosetta/channel"
)

type IteratorOption func(channel <-chan streams.Document, done chan struct{}) <-chan streams.Document

func WithLimit(depth int) IteratorOption {
	return func(ch <-chan streams.Document, done chan struct{}) <-chan streams.Document {
		return channel.Limit(depth, ch, done)
	}
}
