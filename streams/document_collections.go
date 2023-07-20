package streams

import (
	"github.com/benpate/hannibal/vocab"
)

func (document Document) Inbox() Document {
	return document.Get(vocab.PropertyInbox).Collection()
}

func (document Document) Outbox() Document {
	return document.Get(vocab.PropertyOutbox).Collection()
}

func (document Document) Following() Document {
	return document.Get(vocab.PropertyFollowing).Collection()
}

func (document Document) Followers() Document {
	return document.Get(vocab.PropertyFollowers).Collection()
}

func (document Document) Liked() Document {
	return document.Get(vocab.PropertyLiked).Collection()
}

func (document Document) Blocked() Document {
	return document.Get(vocab.PropertyBlocked).Collection()
}
