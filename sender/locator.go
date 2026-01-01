package sender

import "iter"

// Locator defines a service that can locate ActivityPub
// actors and collections based on their URLs
type Locator interface {

	// Actor retrieves a single Actor by its URL.
	Actor(url string) (Actor, error)

	// Recipients a RangeFunc iterator that containing the
	// inbox URLs for every actor that is addressed by
	// the provided URLs.
	// This method should look up individual actors, as well
	// as collections (such as Followers, Circles, etc.)
	// This method is not expected to de-duplicate inbox adddresses;
	// this action will be performed by the Outbox itself.
	Recipient(url string) (iter.Seq[string], error)
}
