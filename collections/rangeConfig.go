package collections

// RangeConfig collects the optional parameters that shape a collection traversal.
type RangeConfig struct {
	maxPages     int
	maxDocuments int
}

// newRangeConfig builds a RangeConfig with default values, then applies the provided options.
func newRangeConfig(options ...Option) RangeConfig {

	// Start with default values
	config := RangeConfig{
		maxPages:     defaultMaxPages,
		maxDocuments: defaultMaxDocuments,
	}

	// Apply each option to the config
	for _, option := range options {
		option(&config)
	}

	// Choices were made.
	return config
}
