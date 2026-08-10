package datetime

import (
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// as2DateTime matches the RFC3339 "date-time" production required by
// AS2 Core section 2.3. https://www.w3.org/TR/activitystreams-core/#dates
var as2DateTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(Z|[+-]\d{2}:\d{2})$`)

// TestFormat confirms Format emits RFC3339, normalized to UTC.
func TestFormat(t *testing.T) {

	// A fixed instant expressed in a non-UTC zone must format as its UTC equivalent.
	instant := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*60*60))

	require.Equal(t, "2026-01-02T22:04:05Z", Format(instant))
	require.Regexp(t, as2DateTime, Format(instant))
}

// TestFormat_NotAnHTTPDate confirms Format does not emit an HTTP date.
func TestFormat_NotAnHTTPDate(t *testing.T) {

	// This is the regression guard for the defect this package exists to fix:
	// "published" was being emitted in the HTTP "Date" header format, which
	// peers parsing AS2 as ISO 8601 cannot read.

	instant := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	formatted := Format(instant)

	// It parses as RFC3339...
	parsed, err := time.Parse(time.RFC3339, formatted)
	require.NoError(t, err)
	require.True(t, parsed.Equal(instant))

	// ...and specifically does NOT parse as an HTTP date.
	_, err = time.Parse(time.RFC1123, formatted)
	require.Error(t, err)
	require.NotEqual(t, "Fri, 02 Jan 2026 15:04:05 GMT", formatted)
}

// TestFormat_SubSecondTruncated confirms fractional seconds are dropped,
// keeping output to the plain date-time production.
func TestFormat_SubSecondTruncated(t *testing.T) {

	instant := time.Date(2026, time.January, 2, 15, 4, 5, 123456789, time.UTC)

	require.Equal(t, "2026-01-02T15:04:05Z", Format(instant))
}

// TestFormat_Zero confirms a zero time yields an empty string.
func TestFormat_Zero(t *testing.T) {

	// Callers omit the property entirely instead of claiming a publication date
	// at the Unix epoch.

	require.Equal(t, "", Format(time.Time{}))
	require.Equal(t, "", FromUnixMilli(0))
	require.Equal(t, "", FromUnixSeconds(0))
}

// TestFormat_OutOfRangeYear confirms sentinel timestamps do not escape as
// non-conformant values.
func TestFormat_OutOfRangeYear(t *testing.T) {

	// RFC3339 fixes the year at four digits, so a sentinel such as
	// math.MaxInt64 seconds would otherwise render a 12-digit year.

	require.Equal(t, "", FromUnixSeconds(math.MaxInt64))
	require.Equal(t, "", Format(time.Unix(math.MaxInt64, 0)))
	require.Equal(t, "", Format(time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, "", Format(time.Date(-1, time.January, 1, 0, 0, 0, 0, time.UTC)))

	// The boundaries themselves remain representable. Note that year 1 January 1
	// at midnight UTC is Go's zero time, so the low boundary is tested one
	// second later -- at midnight it is correctly empty for being zero, not for
	// being out of range.
	require.Equal(t, "9999-12-31T23:59:59Z", Format(time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)))
	require.Equal(t, "0001-01-01T00:00:01Z", Format(time.Date(1, time.January, 1, 0, 0, 1, 0, time.UTC)))
}

// TestFromUnix confirms the two epoch constructors agree on the same instant.
func TestFromUnix(t *testing.T) {

	instant := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)
	expected := "2026-01-02T15:04:05Z"

	require.Equal(t, expected, FromUnixSeconds(instant.Unix()))
	require.Equal(t, expected, FromUnixMilli(instant.UnixMilli()))

	// Handing milliseconds to the seconds constructor is the mistake these two
	// functions exist to prevent. It yields a year far outside the RFC3339 range.
	require.NotEqual(t, expected, FromUnixSeconds(instant.UnixMilli()))
}

// TestNow confirms the convenience constructor produces a conformant value.
func TestNow(t *testing.T) {

	formatted := Now()

	require.Regexp(t, as2DateTime, formatted)

	parsed, err := time.Parse(time.RFC3339, formatted)
	require.NoError(t, err)
	require.WithinDuration(t, time.Now(), parsed, time.Minute)
}
