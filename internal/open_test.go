package internal

import (
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestOpen covers the default connector, which reports a driver it cannot find.
func TestOpen(t *testing.T) {
	db, err := Open("no-such-driver", "")
	assert.Error(t, err)
	assert.Nil(t, db)
}
