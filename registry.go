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
// Holding a read lock across a module's Start is what deadlocked the process
// when that module called Find or Register.

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

// registration is a registered module: an instance, or the constructor a
// clone calls to make one.
type registration struct {
	module  Module
	factory func() Module
}

// instance returns the module to start, calling the constructor when the
// registration is one.
func (e registration) instance() Module {
	if e.factory != nil {
		return e.factory()
	}
	return e.module
}

// Register adds a Module to the registry.
//
// Deprecated: use RegisterFunc. One value is shared by every platform that
// clones the registry, including the generations of a reload, so its state
// outlives the platform it was started with.
func (r *Registry) Register(m Module) {
	r.register(m)
}

// RegisterFunc adds a module constructor to the registry. Clone calls it,
// so every platform starts a module of its own.
func (r *Registry) RegisterFunc(f func() Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = append(r.registrations, registration{factory: f})
}

// register adds a module value, and is what the platform registers into,
// its registry being one platform's own. It takes the lock: the callers
// outside this file hold none of their own.
func (r *Registry) register(m Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.registrations = append(r.registrations, registration{module: m})
}

// materialize calls the constructors that have not been called yet. Clone
// does it for a platform; a registry started as it stands does it here.
//
// The caller holds the write lock. A constructor is foreign code, so this is
// only safe for a registry nothing else can reach yet.
func (r *Registry) materialize() {
	for i, e := range r.registrations {
		if e.module == nil {
			r.registrations[i] = registration{module: e.instance()}
		}
	}
}

// Cleanup is sort of a testing.T.Cleanup but for the registry.
// The cleanups are initialized in Start, and ran in Close.
func (r *Registry) Cleanup(fn func(context.Context)) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanup(fn)
}

// cleanup appends a cleanup. The caller holds the write lock.
func (r *Registry) cleanup(fn func(context.Context)) {
	r.cleanups = append(r.cleanups, fn)
}

// Find gets a Module from the registry.
// The target argument can be a pointer or an interface. The function returns true
// if a module matching the type or interface was found and assigned to `target`.
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

// Start will invoke all the modules start functions sequentially.
// If an error occurs, execution is halted and an error is returned.
// The context is passed along for observability and access to the platform.
// The registry's own output goes to the logger of the platform in the
// context; without one it is discarded.
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

// filter will provide a set of enabled modules based on options
// it prints which modules are enabled/disabled to the log.
//
// It takes the registrations rather than reading them off the registry,
// because Name is the module's own code and must not run under the lock.
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

// mount attaches the middleware and then every module's routes. It holds no
// lock: Mount is the module's code and may reach back into the registry.
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

// Close will invoke all the modules close functions in parallel.
// When finished, it will clear the registered modules list, as
// well as any defined middleware and invoked cleanups.
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

// close runs the cleanups in parallel and waits for them. It holds no lock: a
// cleanup is a module's Stop, which may reach back into the registry.
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
