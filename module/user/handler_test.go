package user

import (
	"testing"

	"github.com/titpetric/platform/internal/require"
)

func TestHandler(t *testing.T) {
	h := NewHandler()

	require.NotNil(t, h)
}
