package sigs

import (
	"net/http"
	"time"
)

// dateHeader returns the provided time value in the IMF-fixdate format
// required by the HTTP "Date" header.
// https://www.rfc-editor.org/rfc/rfc9110#http.date
func dateHeader(value time.Time) string {

	// The "Date" header is covered by the signature, so this value must remain
	// readable by parseDateHeader on the receiving end. Normalizing to UTC keeps
	// the output stable regardless of the host's local zone.
	return value.UTC().Format(http.TimeFormat)
}

// parseDateHeader parses the value of an HTTP "Date" header.
func parseDateHeader(value string) (time.Time, error) {

	// Parsing is strict: the "Date" header is covered by the signature, so a
	// value we cannot read is a value we cannot verify.
	return time.Parse(http.TimeFormat, value)
}
