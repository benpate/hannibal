package hannibal

import (
	"net/http"
	"strings"
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

	// If multiple content types are provided, then only check the first one.
	contentType = list.First(contentType, ',')

	// Strip off any parameters from the content type (like charsets and json-ld profiles)
	contentType = list.First(contentType, ';')

	// Remove whitespace around the actual value
	contentType = strings.TrimSpace(contentType)

	// If what remains matches any of these values, then Success!
	switch contentType {
	case vocab.ContentTypeActivityPub,
		vocab.ContentTypeJSON,
		vocab.ContentTypeJSONLD:
		return true
	}

	// Failure.
	return false
}
