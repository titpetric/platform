package ulid_test

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
	"github.com/titpetric/platform/pkg/ulid"
)

func TestULID(t *testing.T) {
	require.Equal(t, 26, len(ulid.String()))
	require.True(t, ulid.Valid("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
	require.False(t, ulid.Valid("124235235"))
}
