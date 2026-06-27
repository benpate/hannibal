package streams

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/uri"
)

// https://www.w3.org/TR/activitypub/#x4-1-actor-objects

// Inbox returns the document's Inbox property.
// https://www.w3.org/TR/activitypub/#inbox
func (document Document) Inbox() Document {
	return document.Get(vocab.PropertyInbox)
}

// Outbox returns the document's Outbox property.
// https://www.w3.org/TR/activitypub/#outbox
func (document Document) Outbox() Document {
	return document.Get(vocab.PropertyOutbox)
}

// Following returns the document's Following property.
// https://www.w3.org/TR/activitypub/#following
func (document Document) Following() Document {
	return document.Get(vocab.PropertyFollowing)
}

// Followers returns the document's Followers property.
// https://www.w3.org/TR/activitypub/#followers
func (document Document) Followers() Document {
	return document.Get(vocab.PropertyFollowers)
}

// Liked returns the document's Liked property.
// https://www.w3.org/TR/activitypub/#liked
func (document Document) Liked() Document {
	return document.Get(vocab.PropertyLiked)
}

// Likes returns the document's Likes property.
// https://www.w3.org/TR/activitypub/#likes
func (document Document) Likes() Document {
	return document.Get(vocab.PropertyLikes)
}

// Blocked returns the document's Blocked property.
// http://w3id.org/fep/c648
func (document Document) Blocked() Document {
	return document.Get(vocab.PropertyBlocked)
}

// Streams returns the document's Streams property.
// https://www.w3.org/TR/activitypub/#streams-property
func (document Document) Streams() Document {
	return document.Get(vocab.PropertyStreams)
}

// PreferredUsername returns the document's PreferredUsername property.
// https://www.w3.org/TR/activitypub/#preferredUsername
func (document Document) PreferredUsername() string {
	return document.Get(vocab.PropertyPreferredUsername).String()
}

// Username returns the document's Username property.
// Alias for https://www.w3.org/TR/activitypub/#preferredUsername
func (document Document) Username() string {
	return document.PreferredUsername()
}

// Featured returns the document's Featured property.
// https://docs.joinmastodon.org/spec/activitypub/#featured
func (document Document) Featured() Document {
	return document.Get(vocab.PropertyFeatured)
}

// UsernameOrID returns the username of the document, if it exists, or the ID of the document if it does not.
func (document Document) UsernameOrID() string {
	if username := document.PreferredUsername(); username != "" {
		return "@" + username + "@" + uri.Hostname(document.ID())
	}
	return document.ID()
}

// Endpoints returns the document's Endpoints property.
// https://www.w3.org/TR/activitypub/#endpoints
func (document Document) Endpoints() Document {
	return document.Get(vocab.PropertyEndpoints)
}
