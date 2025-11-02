package internal

import (
	"testing"

	"github.com/titpetric/platform/internal/require"
)

func TestULID(t *testing.T) {
	require.Equal(t, 26, len(ULID()))
}
