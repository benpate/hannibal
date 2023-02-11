package pub

// isActivityPubContentType returns true if the specified content-type is an ActivityPub content-type
func isActivityPubContentType(contentType string) bool {

	if contentType == ContentTypeActivityPub {
		return true
	}

	if contentType == ContentTypeJSONLD {
		return true
	}

	return false
}
