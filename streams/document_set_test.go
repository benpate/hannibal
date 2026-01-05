package streams

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestSetProperty(t *testing.T) {
	document := NewDocument(mapof.NewAny())
	document.SetProperty("test-property", "test-value")
	require.Equal(t, "test-value", document.Get("test-property").Value())
}

func TestSetID(t *testing.T) {
	document := NewDocument(mapof.NewAny())
	success := document.SetID("test-id")
	require.True(t, success)
	require.Equal(t, "test-id", document.ID())
}

func TestSetID_Nested(t *testing.T) {

	value := mapof.Any{
		"id": "original-id",
		"object": mapof.Any{
			"id": "original-id",
		},
	}
	document := NewDocument(value)
	success := document.SetID("new-id")

	document.Object().SetID("new-id")
	require.True(t, success)
	require.Equal(t, "new-id", document.ID())
	require.Equal(t, "new-id", document.Object().ID())
}
