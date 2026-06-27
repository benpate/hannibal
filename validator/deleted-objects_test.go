package validator

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// TestDeletedObject_NonDelete confirms the validator abstains (Unknown) for any
// activity that is not a Delete.
func TestDeletedObject_NonDelete(t *testing.T) {

	v := NewDeletedObject()

	activity := streams.NewDocument(map[string]any{
		vocab.PropertyType: vocab.ActivityTypeCreate,
	})

	assert.Equal(t, ResultUnknown, v.Validate(blankRequest(), &activity))
}

// TestDeletedObject_MissingObjectID confirms a Delete activity with no object ID
// is rejected as Invalid (there is nothing to confirm as deleted).
//
// Note: the network branches (a 404/410 response proving the object is gone, or a
// 200 proving it still exists) are not exercised here. DeletedObject calls
// remote.Get directly without enabling AllowPrivateIPs, so it cannot reach a
// loopback test server, and the test suite must not make real external requests.
func TestDeletedObject_MissingObjectID(t *testing.T) {

	v := NewDeletedObject()

	activity := streams.NewDocument(map[string]any{
		vocab.PropertyType: vocab.ActivityTypeDelete,
		// No "object" property, so Object().ID() is empty.
	})

	assert.Equal(t, ResultInvalid, v.Validate(blankRequest(), &activity))
}
