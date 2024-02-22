package inbox

type Config struct {
	debugLevel DebugLevel
}

func NewConfig() Config {
	return Config{
		debugLevel: DebugLevelNone,
	}
}

// DebugTerse returns TRUE if this configuration allows Terse-level debugging
func (config Config) DebugTerse() bool {
	return config.debugLevel >= DebugLevelTerse
}

// DebugVerbose returns TRUE if this configuration allows Verbose-level debugging
func (config Config) DebugVerbose() bool {
	return config.debugLevel >= DebugLevelVerbose
}
