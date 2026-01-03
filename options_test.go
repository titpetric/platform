package platform_test

import (
	"testing"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/assert"
)

func TestNewOptions(t *testing.T) {
	t.Setenv("PLATFORM_MODULES", "user,blog")

	opt := platform.NewOptions()
	assert.NotNil(t, opt)

	assert.Equal(t, opt.Modules, []string{"user", "blog"})
}
