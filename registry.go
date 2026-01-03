package platform

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"slices"
	"sync"

	"github.com/titpetric/platform/pkg/telemetry"
)

// Registry provides a programmatic API to manage middleware and modules.
// A module registers middleware and has a contract to enforce lifecycle.
type Registry struct {
	mu sync.RWMutex

	// Modules hold a list of all modules registered. This list
	// is filtered to start/stop only the modules that are enabled.
	// Interacting with a module is subject to concurrency concerns.
	modules    []Module
	middleware []Middleware

	// On registry start when modules start, a cleanup service per module
	// will be registered via this value. On registry close, the slice
	// is cleared. The functions receive the shutdown context.
	cleanups []func(context.Context)
}

// Register adds a Module to the registry.
func (r *Registry) Register(m Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.modules = append(r.modules, m)
}

// Cleanup is sort of a testing.T.Cleanup but for the registry.
// The cleanups are initialized in Start, and ran in Close.
func (r *Registry) Cleanup(fn func(context.Context)) {
	r.cleanups = append(r.cleanups, fn)
}

// Find gets a Module from the registry.
// The target argument can be a pointer or an interface. The function returns true
// if a module matching the type or interface was found and assigned to `target`.
func (r *Registry) Find(target any) bool {
	// target must be a pointer so we can set its underlying value
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		return false
	}
	targetElemType := targetVal.Elem().Type()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, mod := range r.modules {
		moduleVal := reflect.ValueOf(mod)
		moduleType := moduleVal.Type()

		// Direct assignable (module value can be assigned to the target element)
		if moduleType.AssignableTo(targetElemType) {
			targetVal.Elem().Set(moduleVal)
			return true
		}

		// If target is an interface type, check if module implements it.
		// (AssignableTo above already covers the case where targetElemType is
		// the same concrete type; this handles interface implementations.)
		if targetElemType.Kind() == reflect.Interface && moduleType.Implements(targetElemType) {
			targetVal.Elem().Set(moduleVal)
			return true
		}
	}
	return false
}

// Use adds a Middleware to the registry.
func (r *Registry) Use(f Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.middleware = append(r.middleware, f)
}

// Start will invoke all the modules start functions sequentially.
// If an error occurs, execution is halted and an error is returned.
// The context is passed along for observability and access to the platform.
func (r *Registry) Start(ctx context.Context, mux Router, opts *Options) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ctx, span := telemetry.Start(ctx, "registry.Start")
	defer span.End()

	modules, err := r.filter(opts)
	if err != nil {
		return err
	}

	if err := r.start(ctx, modules, opts.Quiet); err != nil {
		return err
	}

	if err := r.mount(ctx, mux, modules, opts.Quiet); err != nil {
		return err
	}

	return nil
}

// filter will provide a set of enabled modules based on options
// it prints which modules are enabled/disabled to the log
func (r *Registry) filter(opts *Options) ([]Module, error) {
	var enabled []Module
	var disabled []string

	for _, mod := range r.modules {
		name := mod.Name()

		// a lifecycle test (main_test.go) catches this at test time
		if name == "" {
			return nil, fmt.Errorf("module %T doesn't return name", mod)
		}

		if len(opts.Modules) > 0 && !slices.Contains(opts.Modules, name) {
			disabled = append(disabled, name)
			continue
		}
		enabled = append(enabled, mod)

	}

	if !opts.Quiet && len(disabled) > 0 {
		log.Printf("[platform] disabled %d modules: %v", len(disabled), disabled)
	}

	return enabled, nil
}

func (r *Registry) mount(ctx context.Context, mux Router, modules []Module, quiet bool) error {
	ctx, span := telemetry.Start(ctx, "registry.mount")
	defer span.End()

	for _, mw := range r.middleware {
		mux.Use(mw)
	}

	for _, mod := range modules {
		if err := mod.Mount(ctx, mux); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) start(ctx context.Context, modules []Module, quiet bool) error {
	ctx, span := telemetry.Start(ctx, "registry.start")
	defer span.End()

	started := make([]string, 0, len(modules))
	for _, mod := range modules {
		name := mod.Name()
		if err := r.startModule(ctx, mod); err != nil {
			return fmt.Errorf("error starting module %s: %w", name, err)
		}
		started = append(started, name)
	}

	if !quiet {
		log.Printf("[platform] started %d modules: %v", len(started), started)
	}

	return nil
}

func (r *Registry) startModule(ctx context.Context, mod Module) error {
	ctx, span := telemetry.Start(ctx, "module.start: "+mod.Name())
	defer span.End()

	r.Cleanup(func(ctx context.Context) {
		r.stopModule(ctx, mod)
	})

	return mod.Start(ctx)
}

func (r *Registry) stopModule(ctx context.Context, mod Module) {
	ctx, span := telemetry.Start(ctx, "module.stop: "+mod.Name())
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			telemetry.CaptureError(ctx, fmt.Errorf("recovered panic: %v", r))
		}
	}()

	if err := mod.Stop(ctx); err != nil {
		telemetry.CaptureError(ctx, err)
	}
}

// Close will invoke all the modules close functions in parallel.
// When finished, it will clear the registered modules list, as
// well as any defined middleware and invoked cleanups.
func (r *Registry) Close(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.close(ctx)

	r.modules = r.modules[:0]
	r.middleware = r.middleware[:0]
	r.cleanups = r.cleanups[:0]
}

func (r *Registry) close(ctx context.Context) {
	ctx, span := telemetry.Start(ctx, "registry.close")
	defer span.End()

	if len(r.cleanups) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(r.cleanups))

		for _, fn := range r.cleanups {
			go func() {
				defer wg.Done()
				fn(ctx)
			}()
		}
		wg.Wait()
	}
}

// Clone provides a copy of the registry for use in the platform.
func (r *Registry) Clone() *Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := &Registry{
		modules:    make([]Module, len(r.modules)),
		middleware: make([]Middleware, len(r.middleware)),
	}

	copy(clone.modules, r.modules)
	copy(clone.middleware, r.middleware)

	return clone
}

// Stats returns counts for modules and middlewares in the registry.
func (r *Registry) Stats() (modules, middleware int) {
	return len(r.modules), len(r.middleware)
}
