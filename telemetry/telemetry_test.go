package telemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTelemetry(t *testing.T) {
	_, span := Start(t.Context(), "test.telemetry")
	defer span.End()

	require.True(t, true)
}
