package streams

// Statistics contains totals for various interactions with a document
type Statistics struct {
	Announces int64 `json:"announces,omitempty" bson:"announces,omitempty"` // Announces is the number of times this document has been announced / reposted
	Replies   int64 `json:"replies,omitempty"   bson:"replies,omitempty"`   // Replies is the number of replies to this document
	Likes     int64 `json:"likes,omitempty"     bson:"likes,omitempty"`     // Likes is the number of times this document has been liked
	Dislikes  int64 `json:"dislikes,omitempty"  bson:"dislikes,omitempty"`  // Dislikes is the number of times this document has been disliked
}

// NewStatistics returns a fully initialized Statistics object
func NewStatistics() Statistics {
	return Statistics{}
}

func (stats Statistics) IsEmpty() bool {
	return stats.Announces+stats.Replies+stats.Likes+stats.Dislikes == 0
}

func (stats Statistics) NotEmpty() bool {
	return !stats.IsEmpty()
}
