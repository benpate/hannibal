package streams

// Statistics contains totals for various interactions with a document
type Statistics struct {
	Replies   int64 `json:"replies"   bson:"replies,omitempty"`   // Replies is the number of replies to this document
	Likes     int64 `json:"likes"     bson:"likes,omitempty"`     // Likes is the number of times this document has been liked
	Dislikes  int64 `json:"dislikes"  bson:"dislikes,omitempty"`  // Dislikes is the number of times this document has been disliked
	Announces int64 `json:"announces" bson:"announces,omitempty"` // Announces is the number of times this document has been announced / reposted
}

// NewStatistics returns a fully initialized Statistics object
func NewStatistics() Statistics {
	return Statistics{}
}

func (stats Statistics) IsEmpty() bool {
	return !stats.NotEmpty()
}

func (stats Statistics) NotEmpty() bool {
	return stats.HasReplies() || stats.HasLikes() || stats.HasAnnounces() || stats.HasDislikes()
}

func (stats Statistics) HasReplies() bool {
	return stats.Replies > 0
}

func (stats Statistics) HasLikes() bool {
	return stats.Likes > 0
}

func (stats Statistics) HasDislikes() bool {
	return stats.Dislikes > 0
}

func (stats Statistics) HasAnnounces() bool {
	return stats.Announces > 0
}
