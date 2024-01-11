package streams

import (
	"github.com/benpate/domain"
	"github.com/benpate/hannibal/vocab"
)

// https://www.w3.org/TR/activitypub/#x4-1-actor-objects

// https://www.w3.org/TR/activitypub/#inbox
func (document Document) Inbox() Document {
	return document.Get(vocab.PropertyInbox)
}

// https://www.w3.org/TR/activitypub/#outbox
func (document Document) Outbox() Document {
	return document.Get(vocab.PropertyOutbox)
}

// https://www.w3.org/TR/activitypub/#following
func (document Document) Following() Document {
	return document.Get(vocab.PropertyFollowing)
}

// https://www.w3.org/TR/activitypub/#followers
func (document Document) Followers() Document {
	return document.Get(vocab.PropertyFollowers)
}

// https://www.w3.org/TR/activitypub/#liked
func (document Document) Liked() Document {
	return document.Get(vocab.PropertyLiked)
}

// https://www.w3.org/TR/activitypub/#likes
func (document Document) Likes() Document {
	return document.Get(vocab.PropertyLikes)
}

func (document Document) Blocked() Document {
	return document.Get(vocab.PropertyBlocked)
}

func (document Document) Streams() Document {
	return document.Get(vocab.PropertyStreams)
}

func (document Document) PreferredUsername() string {
	return document.Get(vocab.PropertyPreferredUsername).String()
}

func (document Document) Username() string {
	return document.PreferredUsername()
}

func (document Document) UsernameOrID() string {
	if username := document.PreferredUsername(); username != "" {
		return "@" + username + "@" + domain.NameOnly(document.ID())
	}
	return document.ID()
}

func (document Document) Endpoints() Document {
	return document.Get(vocab.PropertyEndpoints)
}
