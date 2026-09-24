package clients

import (
	"github.com/benpate/hannibal/streams"
	"golang.org/x/sync/singleflight"
)

// Carpool merges concurrent Loads of the same URL, so that one fetch serves every caller waiting on it.
// Create one per process and share it between client stacks; it must not be copied after first use.
type Carpool struct {
	group singleflight.Group
}

// NewCarpool returns an empty Carpool, ready to share between client stacks
func NewCarpool() *Carpool {
	return &Carpool{}
}

// Client wraps innerClient in a CarpoolClient whose Loads are merged through this Carpool with
// other Loads made by the same signer.
func (carpool *Carpool) Client(innerClient streams.Client, signer string) *CarpoolClient {
	return &CarpoolClient{
		carpool:     carpool,
		innerClient: innerClient,
		signer:      signer,
	}
}
