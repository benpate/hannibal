package hannibal

import (
	"net/http"
	"time"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/list"
)

// TimeFormat returns a string representation of the provided time value,
// using the format designated by the W3C spec: https://www.w3.org/TR/activitystreams-core/#dates
func TimeFormat(value time.Time) string {
	return value.UTC().Format(http.TimeFormat)
}

// IsActivityPubContentType returns TRUE if the provided contentType is a valid ActivityPub content type.
// https://www.w3.org/TR/activitystreams-core/#media-type
func IsActivityPubContentType(contentType string) bool {

	// Strip off any parameters from the content type (like charsets and json-ld profiles)
	contentType = list.First(contentType, ';')

	switch contentType {
	case vocab.ContentTypeActivityPub,
		vocab.ContentTypeJSON,
		vocab.ContentTypeJSONLD:
		return true
	}

	return false
}
