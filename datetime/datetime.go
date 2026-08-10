// Package datetime formats timestamps for ActivityStreams 2.0 date-time
// properties, such as "published", "updated", and "startTime". Use it for
// every date-time value that hannibal writes into an ActivityStreams document.
// https://www.w3.org/TR/activitystreams-core/#dates
package datetime

import "time"

// AS2 Core section 2.3: "All properties with date and time values MUST conform
// to the 'date-time' production in [RFC3339] with the one exception that
// seconds MAY be omitted."
//
// This package is deliberately write-only. Strict on write, lenient on read:
// parsing inbound timestamps stays permissive, and is handled by rosetta's
// convert.TimeWithLocale, whose layout list accepts RFC3339 alongside the HTTP
// and RFC822 formats that peers emit in practice.

// Format returns an AS2-conformant date-time for the provided time value,
// or an empty string if the value cannot be rendered conformantly.
func Format(value time.Time) string {

	// A zero time returns "" so that callers omit the property entirely.
	// An absent published property is well-defined in AS2; a timestamp of
	// 1970-01-01 is not -- it is a claim the object was published at the Unix epoch.
	if value.IsZero() {
		return ""
	}

	// Always normalize to UTC, so the result ends in "Z" rather than a numeric
	// offset. Both are legal RFC3339, but pinning to UTC keeps output stable
	// regardless of the host's local zone.
	value = value.UTC()

	// The RFC3339 date-time production fixes the year at four digits, so Go
	// renders anything larger with extra digits -- a sentinel value like
	// math.MaxInt64 seconds becomes "292277026596-12-04T15:30:07Z", which no
	// conformant parser will accept.
	if year := value.Year(); (year < 1) || (year > 9999) {
		return ""
	}

	return value.Format(time.RFC3339)
}

// FromUnixMilli returns an AS2-conformant date-time for a Unix epoch
// timestamp expressed in milliseconds.
func FromUnixMilli(value int64) string {

	// Zero is treated as "no value" instead of the Unix epoch, per Format.
	if value == 0 {
		return ""
	}

	return Format(time.UnixMilli(value))
}

// FromUnixSeconds returns an AS2-conformant date-time for a Unix epoch
// timestamp expressed in seconds.
func FromUnixSeconds(value int64) string {

	// Zero is treated as "no value" instead of the Unix epoch, per Format.
	if value == 0 {
		return ""
	}

	return Format(time.Unix(value, 0))
}

// Now returns the current time as an AS2-conformant date-time.
func Now() string {
	return Format(time.Now())
}
