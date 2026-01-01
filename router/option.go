package router

type Option func(*ReceiveConfig)

func WithValidators(validators ...Validator) Option {
	return func(config *ReceiveConfig) {
		config.Validators = validators
	}
}
