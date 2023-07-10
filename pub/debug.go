package pub

import "strings"

// DebugLevel is a custom enumeration that defines various levels of debug output.
type DebugLevel uint8

// String implements the Stringer interface, and returns a human-readable version of the DebugLevel.
func (debugLevel DebugLevel) String() string {

	switch debugLevel {

	case DebugLevelNone:
		return "None"

	case DebugLevelTerse:
		return "Terse"

	case DebugLevelVerbose:
		return "Verbose"

	default:
		return "Unknown"
	}
}

// DebugLevelNone is the default debug level.  Using this setting, no debug messages will be printed.
const DebugLevelNone DebugLevel = 0

// DebugLevelTerse is a debug level that prints only the most important messages.
const DebugLevelTerse DebugLevel = 100

// DebugLevelVerbose is a debug level that prints all debug messages.
const DebugLevelVerbose DebugLevel = 200

// packageDebugLevel is the package-level setting for printing debug statements.
var packageDebugLevel DebugLevel = DebugLevelNone

func GetDebugLevel() DebugLevel {
	return packageDebugLevel
}

func IsMinDebugLevel(minimum DebugLevel) bool {
	return (packageDebugLevel >= minimum)
}

// SetDebugLevel sets the package-level debug level.  By default, debugging is off, but
// can be enabled by setting the debug level to `pub.DebugLevelTerse` or `pub.DebugLevelVerbose`.
func SetDebugLevel(level DebugLevel) {
	packageDebugLevel = level
}

func SetDebugLevelString(level string) {
	switch strings.ToLower(level) {
	case "verbose":
		packageDebugLevel = DebugLevelVerbose
	case "terse":
		packageDebugLevel = DebugLevelTerse
	default:
		packageDebugLevel = DebugLevelNone
	}
}
