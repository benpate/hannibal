package inbox

import (
	"github.com/rs/zerolog"
)

// canLog is a silly zerolog helper that returns TRUE
// if the provided log level would be allowed
// (based on the global log level).
// This makes it easier to execute expensive code conditionally,
// for instance: marshalling a JSON object for logging.
func canLog(level zerolog.Level) bool {
	return zerolog.GlobalLevel() <= level
}

// canTrace returns TRUE if zerolog is configured to allow Trace logs
// This function is here for completeness.  It may or may not be used
func canTrace() bool {
	return canLog(zerolog.TraceLevel)
}

// canDebug returns TRUE if zerolog is configured to allow Debug logs
// This function is here for completeness.  It may or may not be used
func canDebug() bool {
	return canLog(zerolog.DebugLevel)
}

// canInfo returns TRUE if zerolog is configured to allow Info logs
// This function is here for completeness.  It may or may not be used
func canInfo() bool {
	return canLog(zerolog.InfoLevel)
}
