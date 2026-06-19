package validator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/stretchr/testify/assert"
)

// blankRequest returns a minimal request for validators that do not inspect it.
func blankRequest() *http.Request {
	return httptest.NewRequest(http.MethodPost, "https://example.com/inbox", nil)
}

// TestNone confirms the None validator always returns Valid (the request was
// authenticated by some external means).
func TestNone(t *testing.T) {

	v := NewNone()
	activity := streams.NewDocument(map[string]any{})

	assert.Equal(t, ResultValid, v.Validate(blankRequest(), &activity))
}

// actorDocument builds an activity whose actor is an embedded object with the
// given ID. (A bare-string actor reference would require a remote dereference,
// which the offline test environment cannot perform, so the actor is inlined.)
func actorDocument(actorID string) streams.Document {
	return streams.NewDocument(map[string]any{
		vocab.PropertyActor: map[string]any{
			vocab.PropertyID:   actorID,
			vocab.PropertyType: vocab.ActorTypePerson,
		},
	})
}

// TestMatchActor confirms the actor must match the configured actor ID.
func TestMatchActor(t *testing.T) {

	v := NewMatchActor("https://example.com/users/alice")

	t.Run("matching actor -> valid", func(t *testing.T) {
		activity := actorDocument("https://example.com/users/alice")
		assert.Equal(t, ResultValid, v.Validate(blankRequest(), &activity))
	})

	t.Run("different actor -> invalid", func(t *testing.T) {
		activity := actorDocument("https://example.com/users/eve")
		assert.Equal(t, ResultInvalid, v.Validate(blankRequest(), &activity))
	})

	t.Run("missing actor -> invalid", func(t *testing.T) {
		activity := streams.NewDocument(map[string]any{})
		assert.Equal(t, ResultInvalid, v.Validate(blankRequest(), &activity))
	})
}
