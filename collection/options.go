package collection

// Option is a functional option that can be used to configure a collection
type Option func(config *Config)

// WithSSEEndpoint sets the URL endpoint for Server-Sent Events (SSE) updates related to this collection
func WithSSEEndpoint(endpoint string) Option {
	return func(config *Config) {
		config.SSEEndpoint = endpoint
	}
}

// WithAttributedTo sets the "attributedTo" property for this collection,
// which indicates the entity responsible for the collection's content
func WithAttributedTo(attributedTo string) Option {
	return func(config *Config) {
		config.AttributedTo = attributedTo
	}
}

// WithAudience sets the "audience" property for this collection,
// which indicates the intended audience for the collection's content
func WithAudience(audience string) Option {
	return func(config *Config) {
		config.Audience = audience
	}
}
