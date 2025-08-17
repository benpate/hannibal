package streams

import "github.com/benpate/hannibal/vocab"

// DocumentCategory returns the higher level category for the provided document type: [Activity, Actor, Collection, Object]
func DocumentCategory(documentType string) string {

	if IsActivity(documentType) {
		return vocab.DocumentCategoryActivity
	}

	if IsActor(documentType) {
		return vocab.DocumentCategoryActor
	}

	if IsCollection(documentType) {
		return vocab.DocumentCategoryCollection
	}

	return vocab.DocumentCategoryObject
}

func IsActivity(documentType string) bool {

	switch documentType {

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

func IsActor(documentType string) bool {

	switch documentType {

	case vocab.ActorTypeApplication,
		vocab.ActorTypeGroup,
		vocab.ActorTypePerson,
		vocab.ActorTypeOrganization,
		vocab.ActorTypeService:
		return true
	}

	return false
}

func IsCollection(documentType string) bool {

	switch documentType {

	case vocab.CoreTypeCollection,
		vocab.CoreTypeCollectionPage,
		vocab.CoreTypeOrderedCollection,
		vocab.CoreTypeOrderedCollectionPage:

		return true
	}

	return false
}

func IsObject(documentType string) bool {

	switch documentType {

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
