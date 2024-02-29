package sigs

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCanTrace(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	require.True(t, canTrace())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	require.False(t, canTrace())
}

func TestCanDebug(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	require.False(t, canDebug())
}

func TestCanInfo(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	require.True(t, canInfo())
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	require.False(t, canInfo())
}
