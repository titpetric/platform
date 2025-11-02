package telemetry

import (
	"testing"

	"github.com/titpetric/platform/internal/require"
)

func TestTelemetry(t *testing.T) {
	_, span := Start(t.Context(), "test.telemetry")
	defer span.End()

	require.True(t, true)
}
