package inbox

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/rs/zerolog"
)

// IsActivityPubContentType returns true if the specified content-type is an ActivityPub content-type
func IsActivityPubContentType(contentType string) bool {

	if contentType == vocab.ContentTypeActivityPub {
		return true
	}

	if contentType == vocab.ContentTypeJSONLD {
		return true
	}

	return false
}

// canLog is a silly zerolog helper that returns TRUE
// if the provided log level would be allowed
// (based on the global log level).
// This makes it easier to execute expensive code conditionally,
// for instance: marshalling a JSON object for logging.
func canLog(level zerolog.Level) bool {
	return zerolog.GlobalLevel() <= level
}
