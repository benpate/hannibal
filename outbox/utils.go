package outbox

import (
	"iter"

	"github.com/rs/zerolog"
)

// canDebug returns TRUE if zerolog is configured to allow Debug logs
func canDebug() bool {
	return canLog(zerolog.DebugLevel)
}

// canLog is a silly zerolog helper that returns TRUE
// if the provided log level would be allowed
// (based on the global log level).
// This makes it easier to execute expensive code conditionally,
// for instance: marshalling a JSON object for logging.
func canLog(level zerolog.Level) bool {
	return zerolog.GlobalLevel() <= level
}

// makeIterator returns an iter.Seq that yields the provided values in order.
func makeIterator[T any](values ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}
}
