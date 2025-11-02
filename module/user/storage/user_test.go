package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserStorage(t *testing.T) {
	s := NewUserStorage(nil)

	require.NotNil(t, s)
}
