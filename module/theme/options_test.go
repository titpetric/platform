package theme

import (
	"testing"

	"github.com/titpetric/platform/internal/require"
)

func TestTheme(t *testing.T) {
	require.NotNil(t, NewOptions())
}
