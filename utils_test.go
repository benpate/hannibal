package hannibal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsActivityPubContentType(t *testing.T) {
	require.True(t, IsActivityPubContentType("application/json"))
	require.True(t, IsActivityPubContentType("application/json; everything after the semicolon is ignored"))
	require.True(t, IsActivityPubContentType("application/json; whocares=notme"))
	require.True(t, IsActivityPubContentType("application/activity+json"))
	require.True(t, IsActivityPubContentType("application/activity+json; charset=utf-8"))
	require.True(t, IsActivityPubContentType("application/ld+json"))
	require.True(t, IsActivityPubContentType("application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\""))

	require.False(t, IsActivityPubContentType("literally anything else"))
	require.False(t, IsActivityPubContentType("application/xml"))
	require.False(t, IsActivityPubContentType("application/xml; whocares=notme"))
	require.False(t, IsActivityPubContentType("image/webp"))
}
