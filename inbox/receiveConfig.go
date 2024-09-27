package inbox

import "github.com/benpate/hannibal/validator"

// ReceiveConfig is a configuration object for the `ReceiveRequest` function.
type ReceiveConfig struct {
	Validators []Validator
}

// NewReceiveConfig creates a new ReceiveConfig object with default settings,
// and applies any provided options to override the defaults.
func NewReceiveConfig(options ...Option) ReceiveConfig {

	result := ReceiveConfig{
		Validators: []Validator{

			// TODO: check Object Integrity Proofs?

			// checks HTTP signatures
			validator.NewHTTPSig(),

			// checks if objects have been deleted
			validator.NewDeletedObject(),
		},
	}

	for _, option := range options {
		option(&result)
	}

	return result
}
