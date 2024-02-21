package inbox

import "github.com/benpate/hannibal/vocab"

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
