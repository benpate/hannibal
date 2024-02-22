package inbox

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
