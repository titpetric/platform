// The platform is an extensible modular system for writing HTTP servers.
//
// 1. Provides a global registry for middleware and module registration
// 2. Provides a lifecycle to the modules for graceful shutdown
// 3. Provides a router the modules can attach to
//
// It's advised to use `platform.Register` from `init` functions.
// Similarly, `platform.Use` should be used from `main` or any
// descendant setup functions. Don't use these functions from tests
// as they create a shared state.
//
// It's possible to use the platform in an emperative way.
//
// ```go
// svc := platform.New(platform.NewOptions())
// svg.Use(middleware.Logger)
// svc.Register(user.NewModule())
// ```
//
// The platform lifecycle is extensively tested to ensure no races, no
// goroutine leaks. Each platform object creates a copy of the global
// state and holds scoped allocations only, enabling test parallelism.
package platform

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/pkg/telemetry"
)

// Platform is our world struct.
type Platform struct {
	options *Options

	// server setup
	router   *chi.Mux
	server   *http.Server
	listener net.Listener

	// final shutdown context
	context context.Context
	cancel  context.CancelFunc
	stop    func()
	once    sync.Once

	// registry holds settings for plugins and middleware.
	// It's auto-filled from global scope.
	registry *Registry
}

// New will create a new *Platform object. It is the allocation point
// for each platform instance. If no options are passed, the defaults are in use.
// The defaults options are provided by NewOptions().
func New(options *Options) *Platform {
	if options == nil {
		options = NewOptions()
	}

	p := &Platform{
		options: options,
		router:  chi.NewRouter(),
		stop:    func() {},
	}

	// Set up the default registry.
	p.registry = global.registry.Clone()

	// Set up final shutdown signal.
	p.context, p.cancel = context.WithCancel(context.Background())
	return p
}

// Register will add a registry.Module into the internal platform registry.
// This function should be called before Serve is called.
func (p *Platform) Register(m Module) {
	p.registry.Register(m)
}

// Use will add a middleware to the internal platform registry.
// This function should be called before Serve is called.
func (p *Platform) Use(m Middleware) {
	p.registry.Use(m)
}

// Stats will report how many middlewares and plugins are added to the registry.
func (p *Platform) Stats() (int, int) {
	return p.registry.Stats()
}

// Find fills target with the module matching the type.
func (p *Platform) Find(target any) bool {
	return p.registry.Find(target)
}

// Start will start the server and print the registered routes.
// It respects cancellation from the passed context, as well as
// sets up signal notification to respond to SIGTERM.
func (p *Platform) Start(ctx context.Context) error {
	if err := p.setup(ctx); err != nil {
		return fmt.Errorf("error in platform setup: %w", err)
	}

	// If the program receives a SIGTERM, trigger shutdown.
	sigctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	p.stop = stop

	go func() {
		<-sigctx.Done()
		if !p.options.Quiet {
			log.Println("caught sigterm, stopping server")
		}
		p.Stop()
	}()

	// Start the server.
	go func() {
		if err := p.server.Serve(p.listener); err != nil && err != http.ErrServerClosed {
			telemetry.CaptureError(p.context, err)
		}
	}()

	// Print registered routes.
	if !p.options.Quiet {
		internal.PrintRoutes(p.router)
	}

	return nil
}

func (p *Platform) setup(startCtx context.Context) error {
	// set up context for module start
	ctx := platformContext.SetContext(startCtx, p)
	ctx, span := telemetry.Start(ctx, "platform.setup")
	defer span.End()

	if err := p.registry.Start(ctx, p.router, p.options); err != nil {
		return fmt.Errorf("registry: %w", err)
	}

	if err := p.setupListener(); err != nil {
		return fmt.Errorf("setting up listener: %w", err)
	}

	p.server = &http.Server{
		Handler: p.setupRequestContext(p.router),
	}

	return nil
}

// setupRequestContext will bind *Platform to the request context.
func (p *Platform) setupRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		platformContext.Set(r, p)

		next.ServeHTTP(w, r)
	})
}

func (p *Platform) setupListener() error {
	// Set up server listener.
	listener, err := net.Listen("tcp", p.options.ServerAddr)
	if err != nil {
		return err
	}
	p.listener = listener

	if !p.options.Quiet {
		log.Println("Server listening on", p.listener.Addr().String(), p.URL())
	}
	return nil
}

// Context returns the cancellation context for the service.
// When the context finishes, the server has shut down.
func (p *Platform) Context() context.Context {
	return p.context
}

// Wait will pause until the server is shut down.
func (p *Platform) Wait() {
	// Wait for Stop() to be invoked.
	<-p.context.Done()
}

// URL gives the e2e endpoint URL for requests.
func (p *Platform) URL() string {
	listenAddr := p.listener.Addr().String()
	_, port, _ := net.SplitHostPort(listenAddr)
	return "http://127.0.0.1:" + port
}

// Stop will gracefully shutdown the server and then cancel the server context when done.
//
// Stop is an important part of the lifecycle tests. When closing the registry,
// each plugins Stop function gets invoked in parallel. This enables the plugin
// to clear background goroutine event loops, or flush a dirty buffer to storage.
//
// Only after the server has fully shut down does the internal context get cancelled.
func (p *Platform) Stop() {
	p.once.Do(func() {
		// Give a 5 second timeout for a graceful shutdown.
		ctx, cancel := context.WithTimeout(p.Context(), 5*time.Second)
		defer cancel()

		// When done, exit main. It's waiting for the cancelled context there.
		defer func() {
			p.stop()
			p.cancel()
			p.registry.Close(p.context)
		}()

		// Capture error to telemetry sink.
		telemetry.CaptureError(p.context, p.server.Shutdown(ctx))
	})
	return
}
