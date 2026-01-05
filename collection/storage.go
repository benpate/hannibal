package collection

import (
	"iter"

	"github.com/benpate/rosetta/mapof"
)

type IdentifierFunc func() string

type CounterFunc func() (int64, error)

type IteratorFunc func(startIndex string) (iter.Seq[mapof.Any], error)
