package pub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

func GetProfile(remoteID string) (mapof.Any, error) {
	return Get(remoteID)
}

func GetInboxURL(remoteID string) (string, error) {
	profile, err := Get(remoteID)
	return profile.GetString("inbox"), err
}

func GetOutboxURL(remoteID string) (string, error) {
	profile, err := Get(remoteID)
	return profile.GetString("outbox"), err
}

func GetFollowersURL(remoteID string) (string, error) {
	profile, err := Get(remoteID)
	return profile.GetString("followers"), err
}

func GetFollowingURL(remoteID string) (string, error) {
	profile, err := Get(remoteID)
	return profile.GetString("following"), err
}

/******************************************
 * Basic HTTP Operations
 ******************************************/

func Get(remoteID string) (mapof.Any, error) {

	// TODO: Some values should be cached internally in this package

	result := mapof.NewAny()

	transaction := remote.Get(remoteID).
		Header("Accept", "application/activity+json").
		Response(&result, nil)

	if err := transaction.Send(); err != nil {
		return result, derp.Wrap(err, "activitypub.GetProfile", "Error getting profile", remoteID)
	}

	return result, nil
}

// Post sends an ActivityStream to a remote ActivityPub service
// actor: The Actor that is sending the request
// activity: The ActivityStream that is being sent
// targetID: The ID of the Actor that will receive the request
func Post(actor Actor, activity mapof.Any, targetID string) error {

	// Try to get the source profile that we're going to follow
	target, err := GetProfile(targetID)

	if err != nil {
		return derp.Wrap(err, "activitypub.Follow", "Error getting source profile", targetID)
	}

	// Try to get the actor's inbox from the actor ActivityStream.
	// TODO: LOW: Is there a better / more reliable way to do this?
	inbox := target.GetString("inbox")

	if inbox == "" {
		return derp.NewInternalError("activitypub.Follow", "Unable to find 'inbox' in target profile", targetID, target)
	}

	// Send the request to the target Actor's inbox
	transaction := remote.Post(inbox).
		Accept(vocab.ContentTypeActivityPub).
		ContentType(vocab.ContentTypeActivityPub).
		Use(RequestSignature(actor)).
		JSON(activity)

	if err := transaction.Send(); err != nil {
		return derp.Wrap(err, "activitypub.Follow", "Error sending Follow request", inbox)
	}

	// Done!
	return nil
}
