package collection

import (
	"iter"

	"github.com/benpate/rosetta/mapof"
)

// IdentifierFunc returns the unique identifier for a collection item.
type IdentifierFunc func() string

// CounterFunc returns the total number of items in a collection.
type CounterFunc func() (int64, error)

// IteratorFunc iterates over the items in a collection.
type IteratorFunc func(startIndex string) (iter.Seq[mapof.Any], error)
