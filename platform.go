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

	// registry holds settings for plugins and middleware.
	// It's auto-filled from global scope.
	registry *Registry
}

// New is a shorthand for NewPlatform, using default options.
func New() (*Platform, error) {
	return NewPlatform(nil)
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
		stop:    func() {},
	}

	// Set up the default registry.
	p.registry = global.registry.Clone()

	// Set up final shutdown signal.
	p.context, p.cancel = context.WithCancel(context.Background())
	return p, nil
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

// Start will start the server and print the registered routes.
// It respects cancellation from the passed context, as well as
// sets up signal notification to respond to SIGTERM.
func (p *Platform) Start(ctx context.Context) error {
	if err := p.registry.Start(p.router); err != nil {
		return err
	}

	// Set up server listener.
	listener, err := net.Listen("tcp", p.options.ServerAddr)
	if err != nil {
		return err
	}
	p.listener = listener

	if !p.options.Quiet {
		log.Println("Server listening on", p.listener.Addr().String())
	}

	// Set up the server.
	p.server = &http.Server{
		Handler: p.router,
	}

	// If the program receives a SIGTERM, trigger shutdown.
	sigctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	p.stop = stop

	go func() {
		<-sigctx.Done()
		if !p.options.Quiet {
			log.Println("caught sigterm, stopping server")
		}
		if err := p.Stop(); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	// Start the server.
	go func() {
		if err := p.server.Serve(p.listener); err != nil && err != http.ErrServerClosed {
			// Fatalf would skip module cancellation. This just logs the shutdown issue.
			log.Printf("server error: %v", err)
		}
	}()

	// Print registered routes.
	internal.PrintRoutes(p.router)

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
// This is used by `Wait` to signal the service has shut down and the program can cleanly exit.
func (p *Platform) Stop() error {
	// When done, exit main. It's waiting for the cancelled context there.
	defer p.cancel()
	defer p.stop()

	// Clear registry on shutdown.
	defer p.registry.Close()

	// Give a 5 second timeout for a graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := p.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
