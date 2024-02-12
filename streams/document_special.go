package streams

import (
	"strconv"

	"github.com/benpate/hannibal/vocab"
)

/******************************************
 * Type Detection
 ******************************************/

// IsActivity returns TRUE if this document represents an Activity
func (document Document) IsActivity() bool {

	switch document.Type() {

	case vocab.ActivityTypeAccept,
		vocab.ActivityTypeAdd,
		vocab.ActivityTypeAnnounce,
		vocab.ActivityTypeArrive,
		vocab.ActivityTypeBlock,
		vocab.ActivityTypeCreate,
		vocab.ActivityTypeDelete,
		vocab.ActivityTypeDislike,
		vocab.ActivityTypeFlag,
		vocab.ActivityTypeFollow,
		vocab.ActivityTypeIgnore,
		vocab.ActivityTypeInvite,
		vocab.ActivityTypeJoin,
		vocab.ActivityTypeLeave,
		vocab.ActivityTypeLike,
		vocab.ActivityTypeListen,
		vocab.ActivityTypeMove,
		vocab.ActivityTypeOffer,
		vocab.ActivityTypeQuestion,
		vocab.ActivityTypeReject,
		vocab.ActivityTypeRead,
		vocab.ActivityTypeRemove,
		vocab.ActivityTypeTentativeReject,
		vocab.ActivityTypeTentativeAccept,
		vocab.ActivityTypeTravel,
		vocab.ActivityTypeUndo,
		vocab.ActivityTypeUpdate,
		vocab.ActivityTypeView:
		return true
	}

	return false
}

// NotActivity returns TRUE if this document does NOT represent an Activity
func (document Document) NotActivity() bool {
	return !document.IsActivity()
}

// IsActor returns TRUE if this document represents an Actor
func (document Document) IsActor() bool {

	switch document.Type() {

	case vocab.ActorTypeApplication,
		vocab.ActorTypeGroup,
		vocab.ActorTypeOrganization,
		vocab.ActorTypePerson,
		vocab.ActorTypeService:
		return true
	}

	return false
}

// NotActor returns TRUE if this document does NOT represent an Actor
func (document Document) NotActor() bool {
	return !document.IsActor()
}

// IsCollection returns TRUE if this document represents a Collection or CollectionPage
func (document Document) IsCollection() bool {

	switch document.Type() {
	case vocab.CoreTypeCollection,
		vocab.CoreTypeCollectionPage,
		vocab.CoreTypeOrderedCollection,
		vocab.CoreTypeOrderedCollectionPage:

		return true
	}

	return false
}

// NotCollection returns TRUE if the document does NOT represent a Collection or CollectionPage
func (document Document) NotCollection() bool {
	return !document.IsCollection()
}

// IsObject returns TRUE if this document represents an Object type (Article, Note, etc)
func (document Document) IsObject() bool {

	switch document.Type() {

	case vocab.ObjectTypeArticle,
		vocab.ObjectTypeAudio,
		vocab.ObjectTypeDocument,
		vocab.ObjectTypeEvent,
		vocab.ObjectTypeImage,
		vocab.ObjectTypeNote,
		vocab.ObjectTypePage,
		vocab.ObjectTypePlace,
		vocab.ObjectTypeProfile,
		vocab.ObjectTypeRelationship,
		vocab.ObjectTypeTombstone,
		vocab.ObjectTypeVideo:

		return true
	}

	return false
}

// NotObject returns TRUE if this document does NOT represent an Object type (Article, Note, etc)
func (document Document) NotObject() bool {
	return !document.IsObject()
}

// Statistics returns counts for various interactions: Announces, Replies, Likes, and Dislikes
func (document Document) Statistics() Statistics {
	return document.statistics
}

// HasImage returns TRUE if this document has a valid Image property
func (document Document) HasImage() bool {
	return document.Image().NotNil()
}

// HasContent returns TRUE if this document has a valid Content property
func (document Document) HasContent() bool {
	return document.Content() != ""
}

// HasSummary returns TRUE if this document has a valid Summary property
func (document Document) HasSummary() bool {
	return document.Summary() != ""
}

func (document Document) HasDimensions() bool {
	return document.Width() > 0 && document.Height() > 0
}

func (document Document) AspectRatio() string {

	width := document.Width()
	height := document.Height()

	if width == 0 || height == 0 {
		return "auto"
	}

	ratio := float64(width) / float64(height)
	return strconv.FormatFloat(ratio, 'f', -1, 64)
}

// If this document is an activity (create, update, delete, etc), then
// this method returns the activity's Object.  Otherwise, it returns
// the document itself.
func (document Document) UnwrapActivity() Document {

	// If this is an "Activity" type, the dig deeper into the object
	// to find the actual document.
	// This is recursive because it's possible to have a deep tree
	// such as Announce > Create > Document. Looking at you, Lemmy...
	if document.IsActivity() {
		return document.Object().UnwrapActivity()
	}

	return document
}
