package platform

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/registry"
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

	// registry holds settings for plugins and middleware.
	// It's currently auto-filled from the registry package.
	registry *registry.Registry
}

// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr string

	// Quiet turns down the verbosity in the Platform logging code, set to true in tests.
	Quiet bool
}

// NewOptions provides default options for the platform.
func NewOptions() *Options {
	return &Options{
		ServerAddr: ":8080",
	}
}

// NewPlatform will create a new *Platform object. It is the allocation point
// for each platform instance. If no options are passed, the defaults are in use.
// The defaults options are provided by NewOptions().
func NewPlatform(options *Options) (*Platform, error) {
	if options == nil {
		options = NewOptions()
	}

	p := &Platform{
		options: options,
		router:  chi.NewRouter(),
	}

	// Set up and mount registered routes.
	p.registry = registry.Clone()
	p.registry.Mount(p.router)

	// Set up server listener.
	listener, err := net.Listen("tcp", p.options.ServerAddr)
	if err != nil {
		return nil, err
	}
	p.listener = listener

	// Set up server.
	p.server = &http.Server{
		Handler: p.router,
	}

	// Set up final shutdown signal.
	p.context, p.cancel = context.WithCancel(context.Background())
	return p, nil
}

// AddModule will add a registry.Module into the internal platform registry.
// This function should be called before Serve is called.
func (p *Platform) AddModule(m registry.Module) {
	p.registry.AddModule(m)
}

// AddMiddleware will add a middleware to the internal platform registry.
// This function should be called before Serve is called.
func (p *Platform) AddMiddleware(m registry.MiddlewareFunc) {
	p.registry.AddMiddleware(m)
}

// Stats will report how many middlewares and plugins are added to the registry.
func (p *Platform) Stats() (int, int) {
	return p.registry.Stats()
}

// Serve will start the server and print the registered routes.
// It respects cancellation from the passed context, as well as
// sets up signal notification to respond to SIGTERM.
func (p *Platform) Serve(ctx context.Context) {
	// If the program receives a SIGTERM, trigger shutdown.
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		p.Close()
	}()

	// If the passed context is cancelled, trigger shutdown.
	go func() {
		<-ctx.Done()
		p.Close()
	}()

	// Start the server.
	go func() {
		if !p.options.Quiet {
			log.Println("Server listening on", p.listener.Addr().String())
		}

		if err := p.server.Serve(p.listener); err != nil && err != http.ErrServerClosed {
			// Fatalf would skip module cancellation. This just logs the shutdown issue.
			log.Printf("server error: %v", err)
		}
	}()

	// Print registered routes.
	internal.PrintRoutes(p.router)
}

// Wait will block until the server has shut down.
func (p *Platform) Wait() {
	<-p.context.Done()
}

// Context returns the cancellation context for the service.
// When the context finishes, the server has shut down.
func (p *Platform) Context() context.Context {
	return p.context
}

// URL gives the e2e endpoint URL for requests.
func (p *Platform) URL() string {
	listenAddr := p.listener.Addr().String()
	_, port, _ := net.SplitHostPort(listenAddr)
	return "http://127.0.0.1:" + port
}

// Close will gracefully shutdown the server and then cancel
// the server context when done. The registry gets gracefully
// shut down and cleared.
func (p *Platform) Close() {
	if !p.options.Quiet {
		log.Println("shutting down server...")
	}

	// When done, exit main. It's waiting for the cancelled context there.
	defer p.cancel()

	// Clear registry on shutdown.
	defer p.registry.Close()

	// Give a 5 second timeout for a graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %+v", err)
	}

	if !p.options.Quiet {
		log.Println("server shutdown done")
	}
}
