package collection

// Config represents optional configuration parameters for a collection
type Config struct {
	SSEEndpoint string
}

// NewConfig creates a new Config struct with the provided options
func NewConfig(options ...Option) Config {

	config := Config{}

	for _, option := range options {
		option(&config)
	}

	return config
}

// Option is a functional option that can be used to configure a collection
type Option func(config *Config)

// WithSSEEndpoint sets the URL endpoint for Server-Sent Events (SSE) updates related to this collection
func WithSSEEndpoint(endpoint string) Option {
	return func(config *Config) {
		config.SSEEndpoint = endpoint
	}
}
