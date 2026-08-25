// The platform is an extensible modular system for writing HTTP servers.
//
// 1. Provides a global registry for middleware and module registration
// 2. Provides a lifecycle to the modules for graceful shutdown
// 3. Provides a router the modules can attach to
//
// It's advised to use `platform.RegisterFunc` from `init` functions.
// Similarly, `platform.Use` should be used from `main` or any
// descendant setup functions. Don't use these functions from tests
// as they create a shared state.
//
// It's possible to use the platform in an imperative way.
//
// ```go
// svc := platform.New(platform.NewOptions())
// svc.Use(middleware.Logger)
// svc.Register(user.NewModule())
// ```
//
// The platform lifecycle is extensively tested to ensure no races, no
// goroutine leaks. Each platform object creates a copy of the global
// state and holds scoped allocations only, enabling test parallelism.
// Modules are part of that copy when registered with `RegisterFunc`,
// which is called once per platform. A value registered with the
// deprecated `Register` is shared by every platform in the process.
package platform

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/oida"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/pkg/httpcontext"
)

// Platform is our world struct.
type Platform struct {
	// Logger receives the platform's own output. New sets it to
	// slog.Default(), or to a discarding logger when Options.Quiet is set.
	// Assign to it before Start to reuse the platform logger from a test,
	// or to route output into a consumer application's logger.
	Logger Logger

	options *Options

	// server setup
	router   *chi.Mux
	server   *http.Server
	listener net.Listener

	// served is closed when the accept loop has returned. Stop waits on it,
	// so a manager gets the socket back when Stop is done.
	served chan struct{}

	// final shutdown context
	context  context.Context
	cancel   context.CancelFunc
	stop     func()
	once     sync.Once
	stopping atomic.Bool

	// registry holds settings for plugins and middleware.
	// It's auto-filled from global scope.
	registry *Registry

	// telemetry records requests and serves the debug dashboard. It's nil
	// when Options.Telemetry is disabled, or when the recorder failed to
	// build, in which case instrumentation degrades to no-ops.
	telemetry *TelemetryModule
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
		served:  make(chan struct{}),
	}

	// Set up the platform logger. It's set before anything that logs, and
	// stays assignable so a test or a consumer can take the output over.
	p.Logger = slog.Default()
	if options.Quiet {
		p.Logger = discard
	}

	// Set up the default registry.
	p.registry = global.registry.Clone()

	// Record into oida unless telemetry is off, which is how a host that
	// brought its own recorder keeps this one off the router. An invalid
	// config disables recording rather than failing the service: every
	// instrumentation call below is nil safe.
	if options.Telemetry.Enabled {
		if module, err := NewTelemetryModule(options.Telemetry); err == nil {
			p.telemetry = module
			p.Use(module.Middleware)
			p.Register(module)
		} else {
			p.Logger.Error("telemetry disabled", "error", err)
		}
	}

	// Set up final shutdown signal.
	p.context, p.cancel = context.WithCancel(context.Background())
	return p
}

// Register will add a registry.Module into the internal platform registry.
// This function should be called before Serve is called.
func (p *Platform) Register(m Module) {
	p.registry.register(m)
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
	// Read the logger once. The field is exported, and the goroutine below
	// outlives this call, so it gets a value and not a shared field.
	log := p.logger()

	if err := p.setup(ctx); err != nil {
		return fmt.Errorf("error in platform setup: %w", err)
	}

	// If the program receives a SIGTERM, trigger shutdown.
	sigctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	p.stop = stop

	go func() {
		<-sigctx.Done()

		// Stop releases the signal handler, which cancels this context
		// as a delivered signal does. Only one of the two is news.
		if p.stopping.Load() {
			return
		}

		log.Info("caught sigterm, stopping server")
		p.Stop()
	}()

	// Start the server.
	go func() {
		defer close(p.served)

		if err := p.server.Serve(p.listener); err != nil && err != http.ErrServerClosed {
			oida.RecordError(p.context, err)
		}
	}()

	// Print registered routes.
	internal.PrintRoutes(log, p.router)

	return nil
}

func (p *Platform) setup(startCtx context.Context) error {
	// set up context for module start
	ctx := platformContext.SetContext(startCtx, p)
	ctx = optionsContext.SetContext(ctx, p.options)

	// Startup does not arrive over the network, so it gets a trace of its
	// own. Without one the spans below have nothing to record onto and the
	// dashboard shows nothing until the first request.
	return p.observe(ctx, "platform.setup", func(ctx context.Context) error {
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
	})
}

// observe runs fn inside a background trace, so the spans it starts are
// recorded. Without a telemetry module there is nothing to record into and fn
// runs as it is.
func (p *Platform) observe(ctx context.Context, name string, fn func(context.Context) error) error {
	if p.telemetry == nil {
		return fn(ctx)
	}
	return p.telemetry.Tracer().Observe(ctx, name, fn)
}

// setupRequestContext will bind *Platform to the request context.
func (p *Platform) setupRequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		platformContext.Set(r, p)
		optionsContext.Set(r, p.options)

		next.ServeHTTP(w, r)
	})
}

func (p *Platform) setupListener() error {
	// A Manager hands its platform a listener to serve, so the socket
	// survives a reload. Only bind one when there is none.
	if p.listener == nil {
		listener, err := net.Listen("tcp", p.options.ServerAddr)
		if err != nil {
			return err
		}
		p.listener = listener
	}

	p.logger().Info("server listening", "addr", p.listener.Addr().String(), "url", p.URL())
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
	return listenerURL(p.listener)
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
		p.stopping.Store(true)

		// Give a 5 second timeout for a graceful shutdown.
		ctx, cancel := context.WithTimeout(p.Context(), 5*time.Second)
		defer cancel()

		// When done, exit main. It's waiting for the cancelled context there.
		defer func() {
			p.stop()
			p.cancel()
			p.registry.Close(p.context)
		}()

		// Setup can fail before the server is built, and Stop is still
		// the way a caller releases what did get allocated. There is
		// nothing to shut down in that case.
		if p.server == nil {
			return
		}

		// Capture error to telemetry sink.
		oida.RecordError(p.context, p.server.Shutdown(ctx))

		// Shutdown returns once the connections are done, but the accept
		// loop is a goroutine of its own, and nothing may be accepting
		// when a manager hands the socket to the next generation.
		select {
		case <-p.served:
		case <-ctx.Done():
		}
	})
}

type platformKey struct{}

var platformContext = httpcontext.NewValue[*Platform](platformKey{})

// FromRequest returns the *Platform instance attached to the request.
func FromRequest(r *http.Request) *Platform {
	return platformContext.Get(r)
}

// FromContext returns the *Platform instance attached to the context.
func FromContext(ctx context.Context) *Platform {
	return platformContext.GetContext(ctx)
}

// Start is a shorthand to create a new *Platform instance and
// immediately starts the server listener and handles requests.
func Start(ctx context.Context, options *Options) (*Platform, error) {
	svc := New(options)
	if err := svc.Start(ctx); err != nil {
		return nil, err
	}
	return svc, nil
}
