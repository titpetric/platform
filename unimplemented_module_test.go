package platform

import (
	"context"
	"testing"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal/assert"
)

// TestUnimplementedModule covers the defaults an embedded module inherits.
func TestUnimplementedModule(t *testing.T) {
	module := NewUnimplementedModule("embedded")
	var m Module = module

	assert.Equal(t, "embedded", m.Name())
	assert.NoError(t, m.Start(context.Background()))
	assert.NoError(t, m.Stop(context.Background()))
	assert.NoError(t, m.Mount(context.Background(), chi.NewRouter()))
}

// TestUnimplementedModuleFuncs covers the function fields, which replace a default.
func TestUnimplementedModuleFuncs(t *testing.T) {
	started, stopped, mounted := false, false, false

	module := UnimplementedModule{
		NameFn:  func() string { return "configured" },
		StartFn: func(context.Context) error { started = true; return nil },
		StopFn:  func(context.Context) error { stopped = true; return nil },
		MountFn: func(context.Context, Router) error { mounted = true; return nil },
	}

	assert.Equal(t, "configured", module.Name())
	assert.NoError(t, module.Start(context.Background()))
	assert.NoError(t, module.Stop(context.Background()))
	assert.NoError(t, module.Mount(context.Background(), chi.NewRouter()))

	assert.True(t, started)
	assert.True(t, stopped)
	assert.True(t, mounted)
}
