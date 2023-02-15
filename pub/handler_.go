package pub

import "github.com/benpate/hannibal/vocab"

// isActivityPubContentType returns true if the specified content-type is an ActivityPub content-type
func isActivityPubContentType(contentType string) bool {

	if contentType == vocab.ContentTypeActivityPub {
		return true
	}

	if contentType == vocab.ContentTypeJSONLD {
		return true
	}

	return false
}
