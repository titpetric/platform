package platform

import (
	"testing"

	"github.com/titpetric/platform/pkg/assert"
)

func TestNewOptions(t *testing.T) {
	opt := NewOptions()
	assert.NotNil(t, opt)
}
