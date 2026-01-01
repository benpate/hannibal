package sender

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/ranges"
	"github.com/stretchr/testify/require"
)

func TestSender_GetRecipients(t *testing.T) {

	activity := mapof.Any{vocab.PropertyTo: "https://test.actor.social/"}

	sender := New(testLocator{}, nil)
	rangeRecipients, err := sender.getRecipients(activity)

	require.Nil(t, err)

	recipients := ranges.Slice(rangeRecipients)
	require.Equal(t, 1, len(recipients))
	require.Equal(t, "https://testactor.social/inbox", recipients[0])
}

func TestSender_GetFollowers(t *testing.T) {

	activity := mapof.Any{vocab.PropertyTo: "https://test.actor.social/followers"}

	sender := New(testLocator{}, nil)
	rangeRecipients, err := sender.getRecipients(activity)

	require.Nil(t, err)

	recipients := ranges.Slice(rangeRecipients)
	require.Equal(t, 3, len(recipients))
	require.Equal(t, "https://follower1.social/inbox", recipients[0])
	require.Equal(t, "https://follower2.social/inbox", recipients[1])
	require.Equal(t, "https://follower3.social/inbox", recipients[2])
}
