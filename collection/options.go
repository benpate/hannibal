package collection

// Option is a functional option that can be used to configure a collection
type Option func(config *Config)

// WithSSEEndpoint sets the URL endpoint for Server-Sent Events (SSE) updates related to this collection
func WithSSEEndpoint(endpoint string) Option {
	return func(config *Config) {
		config.SSEEndpoint = endpoint
	}
}

func WithAttributedTo(attributedTo string) Option {
	return func(config *Config) {
		config.AttributedTo = attributedTo
	}
}

func WithAudience(audience string) Option {
	return func(config *Config) {
		config.Audience = audience
	}
}
