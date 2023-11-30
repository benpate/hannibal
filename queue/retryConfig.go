package queue

import "time"

// RetryConfig contains configuration options for how a Retry operation should be handled
type RetryConfig struct {
	ReportErrors bool
	MaxAttempts  int
	Backoff      BackoffFunc
}

// RetryOption is a function that modifies a RetryConfig
type RetryOption func(*RetryConfig)

// NewRetryConfig creates a new RetryConfig with default values and
// applies additional optional function parameters as provided
func NewRetryConfig(options ...RetryOption) RetryConfig {

	result := RetryConfig{
		MaxAttempts:  10, // with exponential backoff ~1 hour
		ReportErrors: false,
		Backoff: func(attempt int) time.Duration {
			return time.Duration(2^attempt) * time.Second
		},
	}

	result.With(options...)

	return result
}

// With applies one or more options to a RetryConfig
func (config *RetryConfig) With(options ...RetryOption) {
	for _, option := range options {
		option(config)
	}
}

// WithMaxAttempts is a RetryOption that sets the maximum number of attempts to retry a task
func WithMaxAttempts(maxAttempts int) RetryOption {
	return func(config *RetryConfig) {
		config.MaxAttempts = maxAttempts
	}
}

// WithReportErrors is a RetryOption that sets whether or not errors should be reported to the error log
func WithReportErrors() RetryOption {
	return func(config *RetryConfig) {
		config.ReportErrors = true
	}
}

// WithExponentialBackoff is a RetryOption that sets the backoff function to an exponential curve
func WithExponentialBackoff() RetryOption {
	return func(config *RetryConfig) {
		config.Backoff = func(attempt int) time.Duration {
			return time.Duration(2^attempt) * time.Second
		}
	}
}

// WithLinearBackoff is a RetryOption that sets the backoff function to a linear curve
func WithLinearBackoff() RetryOption {
	return func(config *RetryConfig) {
		config.Backoff = func(attempt int) time.Duration {
			return time.Duration(attempt) * time.Second
		}
	}
}
