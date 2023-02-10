package vocab

// ValidateActivityType validates the ActivityPub "type" of the given activity.
// It returns "UNKNOWN" if the type is not recognized.
func ValidateActivityType(activityType string) string {

	switch activityType {

	case ActivityTypeAccept,
		ActivityTypeAdd,
		ActivityTypeAnnounce,
		ActivityTypeArrive,
		ActivityTypeBlock,
		ActivityTypeCreate,
		ActivityTypeDelete,
		ActivityTypeDislike,
		ActivityTypeFlag,
		ActivityTypeFollow,
		ActivityTypeIgnore,
		ActivityTypeInvite,
		ActivityTypeJoin,
		ActivityTypeLeave,
		ActivityTypeLike,
		ActivityTypeListen,
		ActivityTypeMove,
		ActivityTypeOffer,
		ActivityTypeQuestion,
		ActivityTypeReject,
		ActivityTypeRead,
		ActivityTypeRemove,
		ActivityTypeTentativeReject,
		ActivityTypeTentativeAccept,
		ActivityTypeTravel,
		ActivityTypeUndo,
		ActivityTypeUpdate,
		ActivityTypeView:

		return activityType
	}

	return Unknown
}

// ValidateActorType validates the ActivityPub "type" of the given actor.
// It returns "UNKNOWN" if the type is not recognized.
func ValidateActorType(actorType string) string {

	switch actorType {

	case ActorTypeApplication,
		ActorTypeGroup,
		ActorTypeOrganization,
		ActorTypePerson,
		ActorTypeService:

		return actorType
	}

	return Unknown
}

// ValidateObjectType validates the ActivityPub "type" of the given object or link.
// It returns "UNKNOWN" if the type is not recognized.
func ValidateObjectType(objectType string) string {

	switch objectType {

	case ObjectTypeArticle,
		ObjectTypeAudio,
		ObjectTypeDocument,
		ObjectTypeEvent,
		ObjectTypeImage,
		ObjectTypeNote,
		ObjectTypePage,
		ObjectTypePlace,
		ObjectTypeProfile,
		ObjectTypeRelationship,
		ObjectTypeTombstone,
		ObjectTypeVideo,
		LinkTypeMention:

		return objectType
	}

	return Unknown
}
