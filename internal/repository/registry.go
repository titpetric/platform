package repository

import (
	"log"
	"runtime/debug"
	"sync"
)

// Registry provides a programmatic API to manage middleware and plugins.
// A plugin registers middleware and has a contract to enforce lifecycle.
type Registry struct {
	mu sync.RWMutex

	plugins    []Plugin
	middleware []Middleware
}

// Add will create a new plugin in the registry.
func (r *Registry) Add(name string, hook func(Router), shutdown func()) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plugins = append(r.plugins, NewPlugin(name, hook, shutdown))
}

// AddModule will create a new plugin in the registry via a Module.
func (r *Registry) AddModule(m Module) {
	r.Add(m.Name(), m.Mount, m.Close)
}

// Mount will mount all the plugins registered.
func (r *Registry) Mount(mux Router) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, mw := range r.middleware {
		mux.Use(mw)
	}

	for _, plugin := range r.plugins {
		plugin.Mount(mux)
	}
}

// AddMiddleware adds a middleware to the server.
func (r *Registry) AddMiddleware(f Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.middleware = append(r.middleware, f)
}

// Close will invoke all the plugins close functions in parallel.
// When finished, it will clear the registered plugins list, as
// well as any defined middleware.
func (r *Registry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(r.plugins))

	for _, plugin := range r.plugins {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf(plugin.name+".Close recovered panic: %v\n%s", r, debug.Stack())
				}
			}()
			plugin.Close()
		}()
	}

	wg.Wait()

	r.plugins = r.plugins[:0]
	r.middleware = r.middleware[:0]
}

// Clone provides a copy of the registry for use in the platform.
// Closing the copy leaves the package global state alone.
func (r *Registry) Clone() *Registry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := &Registry{
		plugins:    make([]Plugin, len(r.plugins)),
		middleware: make([]Middleware, len(r.middleware)),
	}

	copy(clone.plugins, r.plugins)
	copy(clone.middleware, r.middleware)

	return clone
}

// Stats returns counts for plugins and middlewares in the registry.
func (r *Registry) Stats() (plugins, middleware int) {
	return len(r.plugins), len(r.middleware)
}
