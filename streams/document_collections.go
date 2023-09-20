package streams

import (
	"github.com/benpate/hannibal/vocab"
)

func (document Document) Inbox() Document {
	return document.Get(vocab.PropertyInbox)
}

func (document Document) Outbox() Document {
	return document.Get(vocab.PropertyOutbox)
}

func (document Document) Following() Document {
	return document.Get(vocab.PropertyFollowing)
}

func (document Document) Followers() Document {
	return document.Get(vocab.PropertyFollowers)
}

func (document Document) Liked() Document {
	return document.Get(vocab.PropertyLiked)
}

func (document Document) Blocked() Document {
	return document.Get(vocab.PropertyBlocked)
}
