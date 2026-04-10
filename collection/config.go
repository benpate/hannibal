package collection

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// Config represents optional configuration parameters for a collection
type Config struct {
	SSEEndpoint  string
	AttributedTo string
	Audience     string
}

// NewConfig creates a new Config struct with the provided options
func NewConfig(options ...Option) Config {

	config := Config{}

	for _, option := range options {
		option(&config)
	}

	return config
}

// Apply updates a result map with the optional values defined in this Config
func (config Config) Apply(result *mapof.Any) {

	if config.AttributedTo != "" {
		result.SetString(vocab.PropertyAttributedTo, config.AttributedTo)
	}

	if config.Audience != "" {
		result.SetString(vocab.PropertyAudience, config.Audience)
	}

	if config.SSEEndpoint != "" {
		result.SetString(vocab.PropertyEventStream, config.SSEEndpoint)
	}
}
