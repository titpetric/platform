package platform_test

import (
	"testing"

	"github.com/titpetric/platform/internal/assert"
	"github.com/titpetric/platform/pkg/ulid"
)

func TestULID(t *testing.T) {
	assert.Equal(t, 26, len(ulid.String()))
	assert.True(t, ulid.Valid("01ARZ3NDEKTSV4RRFFQ69G5FAV"))
	assert.False(t, ulid.Valid("124235235"))
}
