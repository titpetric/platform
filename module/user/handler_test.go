package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	h := NewHandler()

	require.NotNil(t, h)
}
