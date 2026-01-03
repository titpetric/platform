package platform_test

import (
	"testing"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/assert"
)

func TestNewOptions(t *testing.T) {
	t.Run("filled env", func(t *testing.T) {
		t.Setenv("PLATFORM_MODULES", "user,blog")

		opt := platform.NewOptions()
		assert.Equal(t, opt.Modules, []string{"user", "blog"})
	})

	t.Run("empty env", func(t *testing.T) {
		t.Setenv("PLATFORM_MODULES", "")

		opt := platform.NewOptions()
		assert.Empty(t, opt.Modules)
	})
}
