package router

import (
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/hannibal/validator"
)

type Option func(*ReceiveConfig)

func WithValidators(validators ...Validator) Option {
	return func(config *ReceiveConfig) {
		config.Validators = validators
	}
}

// WithPublicKeyFinder configures the HTTP signature validator to use the
// provided public key finder when verifying inbound requests. This replaces
// the default HTTPSig validator (which loads the key from the inbound document)
// while leaving the rest of the validator chain intact.
func WithPublicKeyFinder(keyFinder sigs.PublicKeyFinder) Option {
	return func(config *ReceiveConfig) {
		for index, item := range config.Validators {
			if _, ok := item.(validator.HTTPSig); ok {
				config.Validators[index] = validator.NewHTTPSig(keyFinder)
			}
		}
	}
}
