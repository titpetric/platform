package platform_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal/assert"
)

// reentrantModule reaches back into the registry from its own lifecycle
// methods, which is what docs/modules.md tells a module author to do. Each test
// over it covers the locking convention in registry.go, not the module API.
type reentrantModule struct {
	*platform.UnimplementedModule

	fromStart func(context.Context)
	fromMount func(context.Context)
}

func (m *reentrantModule) Start(ctx context.Context) error {
	if m.fromStart != nil {
		m.fromStart(ctx)
	}
	return nil
}

func (m *reentrantModule) Mount(ctx context.Context, _ platform.Router) error {
	if m.fromMount != nil {
		m.fromMount(ctx)
	}
	return nil
}

// sideModule is what the reentrant module looks for.
type sideModule struct {
	*platform.UnimplementedModule
}

func (m *sideModule) Side() string { return "side" }

type sideAPI interface {
	Side() string
}

// runWithin starts a platform and fails the test if it has not finished within
// the timeout. A registry that locks across module code hangs forever here, so
// the timeout is the assertion.
func runWithin(t *testing.T, timeout time.Duration, mod platform.Module) {
	t.Helper()

	p := platform.New(platform.NewTestOptions())
	p.Register(&sideModule{UnimplementedModule: platform.NewUnimplementedModule("side")})
	p.Register(mod)

	done := make(chan error, 1)
	go func() {
		done <- p.Start(t.Context())
	}()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(timeout):
		t.Fatal("platform did not start: the registry is holding a lock across module code")
	}

	stopped := make(chan struct{})
	go func() {
		p.Stop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(timeout):
		t.Fatal("platform did not stop: the registry is holding a lock across a cleanup")
	}
}

func TestRegistryReentrantFind(t *testing.T) {
	var found bool

	runWithin(t, 5*time.Second, &reentrantModule{
		UnimplementedModule: platform.NewUnimplementedModule("reentrant-find"),
		fromStart: func(ctx context.Context) {
			var api sideAPI
			found = platform.FromContext(ctx).Find(&api)
		},
	})

	assert.True(t, found, "Find from a module Start must resolve a side loaded module")
}

func TestRegistryReentrantRegister(t *testing.T) {
	runWithin(t, 5*time.Second, &reentrantModule{
		UnimplementedModule: platform.NewUnimplementedModule("reentrant-register"),
		fromStart: func(ctx context.Context) {
			platform.FromContext(ctx).Register(&sideModule{
				UnimplementedModule: platform.NewUnimplementedModule("late"),
			})
		},
	})
}

func TestRegistryReentrantUseFromMount(t *testing.T) {
	runWithin(t, 5*time.Second, &reentrantModule{
		UnimplementedModule: platform.NewUnimplementedModule("reentrant-use"),
		fromMount: func(ctx context.Context) {
			platform.FromContext(ctx).Use(func(next http.Handler) http.Handler {
				return next
			})
		},
	})
}

// TestRegistryStatsUnderStop covers Stats reading the slices Close truncates.
// Run with -race.
func TestRegistryStatsUnderStop(t *testing.T) {
	p := platform.New(platform.NewTestOptions())
	p.Register(&sideModule{UnimplementedModule: platform.NewUnimplementedModule("side")})

	assert.NoError(t, p.Start(t.Context()))

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 200 {
			_, _ = p.Stats()
		}
	}()

	p.Stop()
	<-done
}

// TestRegistryCleanupDuringStart covers Cleanup taking the lock now that Start
// does not hold it across module code.
func TestRegistryCleanupDuringStart(t *testing.T) {
	registry := &platform.Registry{}
	registry.RegisterFunc(func() platform.Module {
		return &sideModule{UnimplementedModule: platform.NewUnimplementedModule("side")}
	})

	var ran bool
	registry.Cleanup(func(context.Context) { ran = true })

	assert.NoError(t, registry.Start(t.Context(), chi.NewRouter(), platform.NewTestOptions()))

	registry.Close(t.Context())
	assert.True(t, ran, "a cleanup registered before Start runs on Close")
}
