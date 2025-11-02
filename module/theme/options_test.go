package theme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTheme(t *testing.T) {
	require.NotNil(t, NewOptions())
}
