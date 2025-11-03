package telemetry

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
)

func TestTelemetry(t *testing.T) {
	_, span := Start(t.Context(), "test.telemetry")
	defer span.End()

	require.True(t, true)
}
