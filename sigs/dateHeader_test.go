package sigs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestDateHeader confirms dateHeader renders a time in the IMF-fixdate format,
// normalized to UTC.
func TestDateHeader(t *testing.T) {

	// A fixed instant expressed in a non-UTC zone must format as its UTC equivalent.
	instant := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*60*60))

	require.Equal(t, "Fri, 02 Jan 2026 22:04:05 GMT", dateHeader(instant))
}

// TestDateHeader_RoundTrip confirms the two halves of the "Date" header
// contract agree with each other.
func TestDateHeader_RoundTrip(t *testing.T) {

	// The signature covers the "Date" header, so a value this package writes
	// must be one that it can read back and compare.
	instant := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	parsed, err := parseDateHeader(dateHeader(instant))

	require.NoError(t, err)
	require.True(t, parsed.Equal(instant))
}

// TestParseDateHeader_Invalid confirms values that are not HTTP dates are
// rejected, including the RFC3339 format used for ActivityStreams properties.
func TestParseDateHeader_Invalid(t *testing.T) {

	_, err := parseDateHeader("")
	require.Error(t, err)

	_, err = parseDateHeader("not-a-real-date")
	require.Error(t, err)

	_, err = parseDateHeader("2026-01-02T15:04:05Z")
	require.Error(t, err)
}
