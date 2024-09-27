package validator

import (
	"io"
	"net/http"
	"time"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/rs/zerolog/log"
)

// DeletedObject validates "delete" activities by trying to retrieve the original object.
type DeletedObject struct{}

// NewDeletedObject returns a fully initialized DeletedObject validator.
func NewDeletedObject() DeletedObject {
	return DeletedObject{}
}

// Validate implements the Validator interface, which performs the actual validation.
func (v DeletedObject) Validate(request *http.Request, document *streams.Document) Result {

	const location = "hannibal.validator.DeletedObject"

	// Only validate "Delete" activities
	if document.Type() != vocab.ActivityTypeDelete {
		return ResultUnknown
	}

	// Wait for ten minutes before checking.
	// TODO: This job should be queued (or something) in the future.
	time.Sleep(10 * time.Minute)

	// Retrieve the objectID from the document
	objectID := document.Object().ID()

	if objectID == "" {
		return ResultInvalid
	}

	log.Trace().Str("objectID", objectID).Str("location", location).Msg("Validating DeletedObject")

	// Try to retrieve the original document
	txn := remote.Get(objectID).
		Header("Accept", "application/activity+json")

	if err := txn.Send(); err != nil {

		log.Trace().Err(err).Int("code", derp.ErrorCode(err)).Str("location", location).Msg("Received error code")

		// If the document is marked "gone" or "not found",
		// then this "delete" transaction is valid.
		switch derp.ErrorCode(err) {
		case http.StatusNotFound, http.StatusGone:
			return ResultValid
		case http.StatusMovedPermanently:
			log.Trace().Msg("DeletedObject Validator not correctly following redirect.")
			return ResultUnknown
		}

		// We're not expecting this error, so perhaps there's something else going on here.
		derp.Report(derp.Wrap(err, location, "Error retrieving document, but it is not 'gone' or 'not found'"))
		return ResultUnknown
	}

	log.Trace().Str("location", location).Msg("Delete is invalid / document still exists")
	body, err := io.ReadAll(txn.Response().Body)
	log.Trace().Err(err).Msg(string(body))

	// Fall through means that the document still exists, so the "delete" transaction is invalid.
	return ResultInvalid
}
