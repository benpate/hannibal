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

// If this document is an activity (create, update, delete, etc), then
// this method returns the activity's Object.  Otherwise, it returns
// the document itself.
func (document Document) UnwrapActivity() Document {

	if document.IsActivity() {
		return document.Object()
	}

	return document
}
