package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIif confirms iif returns the true value when the condition holds and the
// false value otherwise, for multiple types.
func TestIif(t *testing.T) {

	assert.Equal(t, "yes", iif(true, "yes", "no"))
	assert.Equal(t, "no", iif(false, "yes", "no"))

	assert.Equal(t, 1, iif(true, 1, 2))
	assert.Equal(t, 2, iif(false, 1, 2))
}
