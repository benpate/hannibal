package hannibal

import (
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/benpate/hannibal/vocab"
)

// IsActivityPubContentType returns TRUE if the provided contentType is a valid ActivityPub content type.
// https://www.w3.org/TR/activitystreams-core/#media-type
func IsActivityPubContentType(contentType string) bool {

	// This reads a SINGLE media type, such as a response's `Content-Type`. An "Accept" header is a
	// weighted list, so pass those to IsActivityPubRequest instead.

	// Strip off any parameters from the content type (like charsets and json-ld profiles)
	contentType, _, _ = strings.Cut(contentType, ";")

	// Remove whitespace around the actual value
	contentType = strings.TrimSpace(strings.ToLower(contentType))

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

// IsActivityPubRequest returns TRUE if the request's "Accept" header asks for ActivityPub.
func IsActivityPubRequest(request *http.Request) bool {
	return acceptsActivityPub(request.Header.Get("Accept"))
}

// NotActivityPubRequest returns TRUE if the request's "Accept" header DOES NOT ask for ActivityPub.
func NotActivityPubRequest(request *http.Request) bool {
	return !IsActivityPubRequest(request)
}

// acceptsActivityPub returns TRUE if the provided "Accept" header value asks for ActivityPub.
func acceptsActivityPub(header string) bool {

	// FALSE is where an absent, empty, or wholly unrecognized header lands: ActivityPub is never the
	// default representation of a URL that also serves something else.
	result := false
	bestQuality := float64(0)

	for _, entry := range parseAcceptHeader(header) {

		isActivityPub, recognized := mediaRangeAcceptsActivityPub(entry.mediaRange)

		if !recognized {
			continue
		}

		// A strict comparison is what gives the client's ordering priority on a tie: a later entry
		// of equal quality never displaces an earlier one.
		if entry.quality > bestQuality {
			result = isActivityPub
			bestQuality = entry.quality
		}
	}

	// The people have spoken.
	return result
}

// mediaRangeAcceptsActivityPub maps a single media range onto the answer it selects, and reports
// whether it named anything this library understands.
func mediaRangeAcceptsActivityPub(mediaRange string) (isActivityPub bool, recognized bool) {

	switch mediaRange {

	// RULE: only these three exact media types select ActivityPub -- no wildcard does. A peer always
	// names the type it wants, so an inferred match would only ever mislabel a vague client.
	case vocab.ContentTypeActivityPub,
		vocab.ContentTypeJSONLD,
		vocab.ContentTypeJSON:
		return true, true

	// Wildcards resolve to "not ActivityPub", but are still recognized so that their q-value
	// competes. That keeps `text/*;q=1.0, application/activity+json;q=0.5` returning FALSE.
	case vocab.ContentTypeHTML,
		"text/*",
		"*/*":
		return false, true
	}

	// Anything else names nothing recognized, and the caller skips it rather than reading it as a
	// "no" that could decide the question on its own.
	return false, false
}

// acceptEntry is one entry of an "Accept" header, in the order the client listed it.
type acceptEntry struct {
	mediaRange string
	quality    float64
}

// parseAcceptHeader splits an "Accept" header into its entries, in the order the client listed them.
func parseAcceptHeader(header string) []acceptEntry {

	result := make([]acceptEntry, 0, strings.Count(header, ",")+1)

	for _, text := range splitAcceptHeader(header) {

		if entry, ok := parseAcceptEntry(text); ok {
			result = append(result, entry)
		}
	}

	return result
}

// parseAcceptEntry parses one "Accept" header entry into a media range and a q-value.
func parseAcceptEntry(text string) (acceptEntry, bool) {

	mediaRange, parameters, err := mime.ParseMediaType(text)

	// A malformed *parameter* still yields a usable media range, so only an empty one is fatal
	if mediaRange == "" {
		return acceptEntry{}, false
	}

	// Parse the client's stated preference, which defaults to the maximum when unstated
	quality := float64(1)

	if err == nil {
		if value, exists := parameters["q"]; exists {

			parsed, parseError := strconv.ParseFloat(value, 64)

			// Neither 0 nor 1 is a defensible guess at what an unparseable q-value meant
			if parseError != nil {
				return acceptEntry{}, false
			}

			quality = min(parsed, 1)
		}
	}

	// RULE: "q=0" means "not acceptable" -- a refusal, rather than a weak preference
	if quality <= 0 {
		return acceptEntry{}, false
	}

	return acceptEntry{mediaRange: mediaRange, quality: quality}, true
}

// splitAcceptHeader separates an "Accept" header's comma-delimited entries.
func splitAcceptHeader(header string) []string {

	// Splitting on every comma would be wrong: a parameter value may be a quoted string, and a
	// quoted string is allowed to contain one.
	result := make([]string, 0, strings.Count(header, ",")+1)
	insideQuotes := false
	start := 0

	for index := 0; index < len(header); index++ {

		switch character := header[index]; {

		case character == '"':
			insideQuotes = !insideQuotes

		// Skipping the escaped character keeps an escaped quote from closing the string early
		case character == '\\' && insideQuotes && index+1 < len(header):
			index++

		case character == ',' && !insideQuotes:
			result = append(result, header[start:index])
			start = index + 1
		}
	}

	// The tail after the last comma is an entry too, even when it is empty.
	return append(result, header[start:])
}

// IsUndoableActivity returns TRUE if the provided activityType
// is one that can be undone (as opposed to an activity that must be "Deleted")
func IsUndoableActivity(activityType string) bool {

	switch activityType {

	case vocab.ActivityTypeAnnounce,
		vocab.ActivityTypeDislike,
		vocab.ActivityTypeFollow,
		vocab.ActivityTypeLike,
		vocab.ActivityTypeBlock:
		return true
	}

	return false
}
