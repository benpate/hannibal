package streams

import "github.com/benpate/hannibal/vocab"

// Metadata contains structured metadata for each document, which is useful for collecting/querying records in a database
type Metadata struct {
	HashedID         string `bson:"hashedId,omitempty"`         // HashedID is a unique identifier for this document, used to prevent duplicate records
	DocumentCategory string `bson:"documentCategory,omitempty"` // High-level category of the document [Activity, Actor, Object, Collection]
	RelationType     string `bson:"relationType,omitempty"`     // If this document is related to another document, this contains the type of relation [Reply, Announce, Like, Dislike]
	RelationHref     string `bson:"relationHref,omitempty"`     // If this document is related to another document, this contains the URL of the related document
	Replies          int64  `bson:"replies,omitempty"`          // Replies is the number of replies to this document
	Announces        int64  `bson:"announces,omitempty"`        // Announces is the number of times this document has been announced / reposted
	Likes            int64  `bson:"likes,omitempty"`            // Likes is the number of times this document has been liked
	Dislikes         int64  `bson:"dislikes,omitempty"`         // Dislikes is the number of times this document has been disliked
}

// NewMetadata returns a fully initialized Metadata object
func NewMetadata() Metadata {
	return Metadata{}
}

// IsActor returns TRUE if this document is one of several "Actor" types [Application, Group, Organization, Person, Service]
func (metadata Metadata) IsActor() bool {
	return metadata.DocumentCategory == vocab.DocumentCategoryActor
}

// IsObject returns TRUE if this document is one of several "Object" types [Image, Video, Audio, Document, and others]
func (metadata Metadata) IsObject() bool {
	return metadata.DocumentCategory == vocab.DocumentCategoryObject
}

// IsCollection returns TRUE if this document is one of several "Collection" types [Collection, CollectionPage, OrderedCollection, OrderedCollectionPage]
func (metadata Metadata) IsCollection() bool {
	return metadata.DocumentCategory == vocab.DocumentCategoryCollection
}

// HasReplies returns TRUE if this document has one or more Replies
func (metadata Metadata) HasReplies() bool {
	return metadata.Replies > 0
}

// HasAnnounces returns TRUE if this document has one or more Announces
func (metadata Metadata) HasAnnounces() bool {
	return metadata.Announces > 0
}

// HasLikes returns TRUE if this document has one or more Likes
func (metadata Metadata) HasLikes() bool {
	return metadata.Likes > 0
}

// HasDislikes returns TRUE if this document has one or more Dislikes
func (metadata Metadata) HasDislikes() bool {
	return metadata.Dislikes > 0
}

// HasRelationship returns TRUE if this document has a relationship
func (metadata Metadata) HasRelationship() bool {
	if metadata.RelationType == "" {
		return false
	}

	if metadata.RelationHref == "" {
		return false
	}

	return true
}

func (metadata *Metadata) SetRelationCount(relationType string, count int64) {

	switch relationType {

	case vocab.RelationTypeReply:
		metadata.Replies = count

	case vocab.RelationTypeAnnounce:
		metadata.Announces = count

	case vocab.RelationTypeLike:
		metadata.Likes = count

	case vocab.RelationTypeDislike:
		metadata.Dislikes = count
	}

}
