package telemetry

import (
	"testing"

	"github.com/titpetric/platform/internal/require"
)

func TestNewMonitor(t *testing.T) {
	m := NewMonitor()

	require.NotNil(t, m)
	require.False(t, m.Enabled())

	m.Touch("test")

	m.SetEnabled(true)
	m.Touch("test")
	m.Touch("test")

	require.True(t, m.Enabled())
}
