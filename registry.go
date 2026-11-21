package platform

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"sync"

	"github.com/titpetric/oida"
)

// Locking convention: an exported method takes r.mu and delegates to an
// unexported one of the same name, which assumes the caller holds it. Find
// and find, Stats and stats, Cleanup and cleanup, Close and close. A method
// documented as holding the lock never calls foreign code.
//
// Module Start, Module Mount, a module cleanup and a RegisterFunc constructor
// are all foreign: they can register, add middleware or Find another module,
// and each of those wants this same mutex. Start, Close and Clone therefore
// snapshot what they need under the lock, release it, and only then call out.

// Registry provides a programmatic API to manage middleware and modules.
// A module registers middleware and has a contract to enforce lifecycle.
type Registry struct {
	mu sync.RWMutex

	// Registrations hold every module registered, in registration order.
	// The list is filtered to start/stop only the modules that are enabled.
	// Interacting with a module is subject to concurrency concerns.
	registrations []registration
	middleware    []Middleware

	// On registry start when modules start, a cleanup service per module
	// will be registered via this value. On registry close, the slice
	// is cleared. The functions receive the shutdown context.
	cleanups []func(context.Context)
}

// registration is a module instance, or the constructor that makes one.
type registration struct {
	module  Module
	factory func() Module
}

// instance returns the module to start, calling the constructor when there is one.
func (e registration) instance() Module {
	if e.factory != nil {
		return e.factory()
	}
	return e.module
}

// Register adds a Module to the registry. Deprecated: use RegisterFunc.
func (r *Registry) Register(m Module) {
	r.register(m)
}

// RegisterFunc adds a module constructor that Clone calls per platform.
func (r *Registry) RegisterFunc(f func() Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = append(r.registrations, registration{factory: f})
}

// register adds a module value, taking the lock itself.
func (r *Registry) register(m Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = append(r.registrations, registration{module: m})
}

// materialize calls the constructors not called yet. The caller holds the write lock.
func (r *Registry) materialize() {
	for i, e := range r.registrations {
		if e.module == nil {
			r.registrations[i] = registration{module: e.instance()}
		}
	}
}

// Cleanup registers a function to run when the registry closes.
func (r *Registry) Cleanup(fn func(context.Context)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanup(fn)
}

// cleanup appends a cleanup. The caller holds the write lock.
func (r *Registry) cleanup(fn func(context.Context)) {
	r.cleanups = append(r.cleanups, fn)
}

// Find fills target, a pointer or an interface, with a matching module.
func (r *Registry) Find(target any) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.find(target)
}

// find gets a Module from the registry. The caller holds the lock.
func (r *Registry) find(target any) bool {
	// target must be a pointer so we can set its underlying value
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Pointer || targetVal.IsNil() {
		return false
	}
	targetElemType := targetVal.Elem().Type()

	for _, e := range r.registrations {
		// A constructor that Clone has not called yet is not an instance
		// to find.
		if e.module == nil {
			continue
		}

		moduleVal := reflect.ValueOf(e.module)
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

// Start starts the modules in order and then mounts them. A failure halts the
// run, and an error ends it.
func (r *Registry) Start(ctx context.Context, mux Router, opts *Options) error {
	ctx, span := oida.Start(ctx, "registry.Start")
	defer span.End()

	log := loggerFromContext(ctx)

	// Snapshot under the lock, then let go of it. Everything below calls into
	// modules, which are free to use the registry.
	r.mu.Lock()
	r.materialize()
	registrations := slices.Clone(r.registrations)
	middleware := slices.Clone(r.middleware)
	r.mu.Unlock()

	modules, err := filter(registrations, log, opts)
	if err != nil {
		return err
	}

	if err := r.start(ctx, modules, log); err != nil {
		return err
	}

	return r.mount(ctx, mux, middleware, modules)
}

// filter returns the modules the options enable, and logs the rest.
func filter(registrations []registration, log Logger, opts *Options) ([]Module, error) {
	var enabled []Module
	var disabled []string

	for _, e := range registrations {
		mod := e.module
		name := mod.Name()

		// a lifecycle test (main_test.go) catches this at test time
		if name == "" {
			return nil, fmt.Errorf("module %T doesn't return name", mod)
		}

		// The allowlist names the application's modules. The recorder
		// the platform registers itself is infrastructure, turned off
		// with Options.Telemetry rather than by omission here.
		if _, builtin := mod.(*TelemetryModule); !builtin {
			if len(opts.Modules) > 0 && !slices.Contains(opts.Modules, name) {
				disabled = append(disabled, name)
				continue
			}
		}
		enabled = append(enabled, mod)

	}

	if len(disabled) > 0 {
		log.Info("modules disabled", "count", len(disabled), "modules", disabled)
	}

	return enabled, nil
}

// mount attaches the middleware and then every module's routes, holding no lock.
func (r *Registry) mount(ctx context.Context, mux Router, middleware []Middleware, modules []Module) error {
	ctx, span := oida.Start(ctx, "registry.mount")
	defer span.End()

	for _, mw := range middleware {
		mux.Use(mw)
	}

	for _, mod := range modules {
		if err := mod.Mount(ctx, mux); err != nil {
			return err
		}
	}

	return nil
}

// start runs every module's Start in registration order. It holds no lock.
func (r *Registry) start(ctx context.Context, modules []Module, log Logger) error {
	ctx, span := oida.Start(ctx, "registry.start")
	defer span.End()

	started := make([]string, 0, len(modules))
	for _, mod := range modules {
		name := mod.Name()
		if err := r.startModule(ctx, mod); err != nil {
			return fmt.Errorf("error starting module %s: %w", name, err)
		}
		started = append(started, name)
	}

	log.Info("modules started", "count", len(started), "modules", started)

	return nil
}

func (r *Registry) startModule(ctx context.Context, mod Module) error {
	ctx, span := oida.Start(ctx, "module.start: "+mod.Name())
	defer span.End()

	r.Cleanup(func(ctx context.Context) {
		r.stopModule(ctx, mod)
	})

	return mod.Start(ctx)
}

func (r *Registry) stopModule(ctx context.Context, mod Module) {
	ctx, span := oida.Start(ctx, "module.stop: "+mod.Name())
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			oida.RecordError(ctx, fmt.Errorf("recovered panic: %v", r))
		}
	}()

	if err := mod.Stop(ctx); err != nil {
		oida.RecordError(ctx, err)
	}
}

// Close runs the cleanups, then clears the modules, middleware and cleanups.
func (r *Registry) Close(ctx context.Context) {
	r.mu.Lock()
	cleanups := r.cleanups

	r.registrations = r.registrations[:0]
	r.middleware = r.middleware[:0]
	// nil rather than truncated: the cleanups above still reference the
	// backing array, and a later append must not write into it.
	r.cleanups = nil
	r.mu.Unlock()

	r.close(ctx, cleanups)
}

// close runs the cleanups in parallel and waits for them, holding no lock.
func (r *Registry) close(ctx context.Context, cleanups []func(context.Context)) {
	ctx, span := oida.Start(ctx, "registry.close")
	defer span.End()

	if len(cleanups) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(cleanups))

	for _, fn := range cleanups {
		go func() {
			defer wg.Done()
			fn(ctx)
		}()
	}
	wg.Wait()
}

// Clone provides a copy of the registry for use in the platform. Modules
// registered as constructors are built here, one set per clone.
func (r *Registry) Clone() *Registry {
	r.mu.RLock()
	registrations := slices.Clone(r.registrations)
	middleware := slices.Clone(r.middleware)
	r.mu.RUnlock()

	clone := &Registry{
		registrations: make([]registration, len(registrations)),
		middleware:    middleware,
	}

	// The constructors are foreign code, so they run with no lock held.
	for i, e := range registrations {
		clone.registrations[i] = registration{module: e.instance()}
	}

	return clone
}

// Stats returns counts for modules and middlewares in the registry.
func (r *Registry) Stats() (modules, middleware int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.stats()
}

// stats counts modules and middleware. The caller holds the lock.
func (r *Registry) stats() (modules, middleware int) {
	return len(r.registrations), len(r.middleware)
}
