package collection

import (
	"iter"

	"github.com/benpate/rosetta/mapof"
)

// Storage interface defines a service that can store
// and retrieve activities (e.g. in a specific User's outbox)
// Services that implement this interface are assumed to already
// contain filtering information such as Actor permissions
// and "after" cursors.
type Storage interface {

	// ID returns the unique ID (URL) of this Storage
	ID() string

	// Count returns the total number of activities in
	// this Storage
	TotalItems() (int, error)

	// Iterator returns a RangeFunc iterator containing
	// all activities in the Storage.
	// Activities MUST be returned in order, and are filtered
	// by the startIndex parameter.
	Iterator(startIndex string) (iter.Seq[mapof.Any], error)
}
