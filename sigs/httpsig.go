// Package sigs implements the IETF draft specification "Signing HTTP Messages"
// https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures

package sigs

import "github.com/rs/zerolog"

func init() {
	// By default, disable all logging.  Applications can override this
	zerolog.SetGlobalLevel(zerolog.Disabled)
}
