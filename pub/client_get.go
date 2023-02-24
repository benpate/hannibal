package pub

import (
	"github.com/benpate/derp"
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
