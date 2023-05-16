package pub

import (
	"encoding/json"
	"fmt"

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

	if packageDebugLevel >= DebugLevelTerse {
		if packageDebugLevel >= DebugLevelVerbose {
			fmt.Println("------------------------------------------")
		}
		fmt.Println("HANNIBAL: Getting RemoteID: " + remoteID)
	}

	transaction := remote.Get(remoteID).
		Header("Accept", "application/activity+json").
		Response(&result, nil)

	if err := transaction.Send(); err != nil {
		return result, derp.Wrap(err, "activitypub.GetProfile", "Error getting profile", remoteID)
	}

	if packageDebugLevel >= DebugLevelVerbose {
		marshalled, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(marshalled))
	}

	return result, nil
}
