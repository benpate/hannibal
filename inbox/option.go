package inbox

type Option func(*Config)

func WithVerbose() Option {
	return func(config *Config) {
		config.debugLevel = DebugLevelVerbose
	}
}

func WithTerse() Option {
	return func(config *Config) {
		config.debugLevel = DebugLevelTerse
	}
}

func WithNoDebug() Option {
	return func(config *Config) {
		config.debugLevel = DebugLevelNone
	}
}
