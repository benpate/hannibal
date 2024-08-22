package hannibal

import (
	"net/http"
	"time"

	"github.com/benpate/hannibal/vocab"
)

// TimeFormat returns a string representation of the provided time value,
// using the format designated by the W3C spec: https://www.w3.org/TR/activitystreams-core/#dates
func TimeFormat(value time.Time) string {
	return value.UTC().Format(http.TimeFormat)
}

// IsActivityPubContentType returns TRUE if the provided contentType is a valid ActivityPub content type.
// https://www.w3.org/TR/activitystreams-core/#media-type
func IsActivityPubContentType(contentType string) bool {

	switch contentType {
	case vocab.ContentTypeActivityPub,
		vocab.ContentTypeJSON,
		vocab.ContentTypeJSONLD,
		vocab.ContentTypeJSONLDWithProfile:
		return true

	}

	return false
}
