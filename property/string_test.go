package property

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {

	value := NewValue("Hello, World!")

	require.Equal(t, "Hello, World!", value.Raw())
	require.Equal(t, "Hello, World!", value.Head().Raw())
	require.Nil(t, value.Tail().Raw())

	require.True(t, IsString(value))
}
