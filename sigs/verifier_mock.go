package sigs

import (
	"net/http"

	"github.com/benpate/derp"
)

// MockVerifier contains all of the settings necessary to verify a request
type MockVerifier struct {
	KeyID   string
	Success bool
}

// NewMockVerifier returns a fully initialized Verifier
func NewMockVerifier(keyID string, success bool) MockVerifier {
	result := MockVerifier{
		KeyID:   keyID,
		Success: success,
	}
	return result
}

// Verify verifies the given http.Request
func (mock *MockVerifier) Verify(request *http.Request, keyFinder PublicKeyFinder) (Signature, error) {

	if mock.Success {
		signature := NewSignature()
		signature.KeyID = mock.KeyID
		return signature, nil
	}

	return NewSignature(), derp.ForbiddenError("hannibal.sigs.MockVerifier.Verify", "MockVerifier is configured to fail")
}
