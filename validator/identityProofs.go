package validator

import (
	"net/http"

	"github.com/benpate/hannibal/streams"
)

// IdentityProof implements FEP-c390 Identity Proofs
// https://codeberg.org/fediverse/fep/src/branch/main/fep/c390/fep-c390.md
type IdentityProof struct{}

func (v IdentityProof) Validate(request *http.Request, document streams.Document) Result {

	return ResultUnknown
}
