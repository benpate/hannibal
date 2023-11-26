package streams

import "github.com/benpate/hannibal/vocab"

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

// If this document is an activity (create, update, delete, etc), then
// this method returns the activity's Object.  Otherwise, it returns
// the document itself.
func (document Document) UnwrapActivity() Document {

	// If this is an "Activity" type, the dig deeper into the object
	// to find the actual document.  This is recursive because it's
	// possible to have a deep tree such as Announce > Create > Document
	// Looking at you, Lemmy...
	if document.IsActivity() {
		return document.Object().UnwrapActivity()
	}

	return document
}
