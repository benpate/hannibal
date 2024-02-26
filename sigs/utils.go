package sigs

import "github.com/rs/zerolog"

func canLog(logLevel zerolog.Level) bool {
	return logLevel >= zerolog.GlobalLevel()
}

func canTrace() bool {
	return canLog(zerolog.TraceLevel)
}

func canDebug() bool {
	return canLog(zerolog.DebugLevel)
}

func canInfo() bool {
	return canLog(zerolog.InfoLevel)
}
