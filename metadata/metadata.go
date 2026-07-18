// Package metadata carries server-computed metadata about ActivityStream documents -- knowledge
// that is ABOUT a document but never part of its wire value. It holds two layers: document facts
// (category, relationships, response counts), which are the same for every viewer and are persisted
// alongside cached documents; and per-viewer moderation Labels, which are attached at load time and
// never persisted or serialized.
package metadata

import "github.com/benpate/hannibal/vocab"

// Metadata contains structured, server-computed metadata about a single document. Document facts
// are shared by every viewer and persisted with the cached document. The Labels result is
// per-viewer, attached at load time and never persisted or serialized.
type Metadata struct {

	// Document facts: identical for every viewer, persisted with the cached document.

	HashedID         string `bson:"hashedId,omitempty"`         // HashedID is a unique identifier for this document, used to prevent duplicate records
	DocumentCategory string `bson:"documentCategory,omitempty"` // High-level category of the document [Activity, Actor, Object, Collection]
	RelationType     string `bson:"relationType,omitempty"`     // If this document is related to another document, this contains the type of relation [Reply, Announce, Like, Dislike]
	RelationHref     string `bson:"relationHref,omitempty"`     // If this document is related to another document, this contains the URL of the related document
	Replies          int64  `bson:"replies,omitempty"`          // Replies is the number of replies to this document
	Announces        int64  `bson:"announces,omitempty"`        // Announces is the number of times this document has been announced / reposted
	Likes            int64  `bson:"likes,omitempty"`            // Likes is the number of times this document has been liked
	Dislikes         int64  `bson:"dislikes,omitempty"`         // Dislikes is the number of times this document has been disliked

	// Labels is the current viewer's moderation verdict for this document. The single bson/json "-"
	// tag keeps EVERYTHING inside it out of shared caches and off the wire, so fields added to
	// Label later are covered automatically -- and a remote server can never spoof it, because the
	// JSON parser fills only the document value, never Metadata.

	Labels LabelSet `bson:"-" json:"-"`
}

// New returns a fully initialized Metadata object.
func New() Metadata {
	return Metadata{}
}

// Clone returns a copy of this Metadata that shares no mutable state with the original.
func (metadata Metadata) Clone() Metadata {
	result := metadata
	result.Labels = metadata.Labels.Clone()
	return result
}

// IsRuleHidden returns TRUE if the current viewer's rules hide this document (a block or a mute).
func (metadata Metadata) IsRuleHidden() bool {
	return metadata.Labels.IsHidden()
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

// SetRelationCount updates the designated relation with a new count,
// returning TRUE if the value has been changed.
func (metadata *Metadata) SetRelationCount(relationType string, count int64) bool {

	switch relationType {

	case vocab.RelationTypeReply:
		if metadata.Replies != count {
			metadata.Replies = count
			return true
		}

	case vocab.RelationTypeAnnounce:
		if metadata.Announces != count {
			metadata.Announces = count
			return true
		}

	case vocab.RelationTypeLike:
		if metadata.Likes != count {
			metadata.Likes = count
			return true
		}

	case vocab.RelationTypeDislike:
		if metadata.Dislikes != count {
			metadata.Dislikes = count
			return true
		}
	}

	return false
}
