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

	t.Run("pidfile from env", func(t *testing.T) {
		t.Setenv("PLATFORM_PIDFILE", "/run/platform.pid")

		opt := platform.NewOptions()
		assert.Equal(t, "/run/platform.pid", opt.PidFile)
	})

	t.Run("no pidfile by default", func(t *testing.T) {
		t.Setenv("PLATFORM_PIDFILE", "")

		opt := platform.NewOptions()
		assert.Empty(t, opt.PidFile)
	})

	// One test binary runs many platforms. If test options picked the
	// variable up, every one of them would write the same path.
	t.Run("test options carry no pidfile", func(t *testing.T) {
		t.Setenv("PLATFORM_PIDFILE", "/run/platform.pid")

		assert.Empty(t, platform.NewTestOptions().PidFile)
	})
}
