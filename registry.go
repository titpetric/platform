package platform

import (
	"log"
	"runtime/debug"
	"sync"
)

// Registry provides a programmatic API to manage middleware and plugins.
// A plugin registers middleware and has a contract to enforce lifecycle.
type Registry struct {
	mu sync.RWMutex

	modules    []Module
	middleware []Middleware
}

func (r *Registry) Register(m Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.modules = append(r.modules, m)
}

// Use adds a middleware to the registry.
func (r *Registry) Use(f Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.middleware = append(r.middleware, f)
}

// Start will invoke all the modules start functions sequentially.
// If an error occurs, execution is halted and an error is returned.
func (r *Registry) Start(mux Router) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, plugin := range r.modules {
		if err := plugin.Start(); err != nil {
			return err
		}
	}

	for _, mw := range r.middleware {
		mux.Use(mw)
	}

	for _, plugin := range r.modules {
		plugin.Mount(mux)
	}

	return nil
}

// Close will invoke all the modules close functions in parallel.
// When finished, it will clear the registered modules list, as
// well as any defined middleware.
func (r *Registry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(r.modules))

	for _, plugin := range r.modules {
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("%s.Close recovered panic: %v\n%s", plugin.Name(), r, debug.Stack())
				}
			}()
			if err := plugin.Stop(); err != nil {
				log.Printf("error in %s: %v", plugin.Name(), err)
			}
		}()
	}

	wg.Wait()

	r.modules = r.modules[:0]
	r.middleware = r.middleware[:0]
}

// Clone provides a copy of the registry for use in the platform.
// Closing the copy leaves the package global state alone.
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
