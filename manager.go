package platform

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// Manager owns the socket, and the platform generation serving on it. It
// replaces the platform under it on SIGHUP, keeping the address it serves.
type Manager struct {
	// Logger receives the manager's own output, and is handed to every
	// platform generation it creates. Defaults as Platform.Logger does.
	Logger Logger

	// Setup runs against every platform generation before it starts.
	// Registration against a platform value belongs here, as a reload
	// discards the value it was made against.
	//
	// Assign it before Start: a reload runs it from the signal handler's
	// goroutine, so assigning it to a running manager is a data race.
	Setup func(*Platform) error

	// Check runs before a reload retires the generation that is serving,
	// and a non-nil error abandons that reload with the generation left
	// alone and still answering requests.
	//
	// It is where reading the configuration a reload would apply belongs.
	// Setup is too late for that: it runs against the new generation, which
	// is built only after the old one has been stopped, so a failure there
	// leaves nothing serving. A signal is not a safe way to hand a process
	// a configuration it has not read yet, and this is what makes it one.
	//
	// Nil, the default, reloads unconditionally.
	//
	// Assign it before Start, as with Setup: the signal handler reads it
	// from a goroutine of its own, so assigning it to a running manager is
	// a data race.
	Check func() error

	options *Options

	// mu serializes the generation swap.
	mu      sync.Mutex
	shared  *sharedListener
	current atomic.Pointer[generation]

	// pid is the file this process records its id in, disabled when
	// Options.PidFile is empty. The manager holds it rather than the
	// platform, because the file records a process and a reload does not
	// make a new one.
	pid pidfile

	// final shutdown context, cancelled when the manager stops
	context context.Context
	cancel  context.CancelFunc
	stop    func()
	once    sync.Once
}

// generation is one platform in the sequence the manager serves.
type generation struct {
	platform *Platform

	// retired marks the teardown as the manager's doing.
	retired atomic.Bool
}

// NewManager creates a manager for the passed options. If no options are
// passed, the defaults from NewOptions() are in use.
func NewManager(options *Options) *Manager {
	if options == nil {
		options = NewOptions()
	}

	m := &Manager{
		options: options,
		stop:    func() {},
		pid:     newPidfile(options.PidFile),
	}

	m.Logger = slog.Default()
	if options.Quiet {
		m.Logger = discard
	}

	m.context, m.cancel = context.WithCancel(context.Background())
	return m
}

// Start binds the listener, starts the first platform generation on it,
// writes Options.PidFile and arms the SIGHUP handler. Cancelling ctx stops
// the manager.
//
// The pidfile is the manager's rather than the generation's: it records a
// process, and a reload replaces the platform without making a new one.
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	listener, err := net.Listen("tcp", m.options.ServerAddr)
	if err != nil {
		return fmt.Errorf("setting up listener: %w", err)
	}
	m.shared = newSharedListener(listener)

	if err := m.startGeneration(ctx); err != nil {
		_ = m.shared.Close()
		return err
	}

	// Once something is serving, and before the handler that a signal would
	// arrive at. A file naming a process that then failed to come up is
	// what a service manager acts on.
	if err := m.pid.write(); err != nil {
		m.retire()
		_ = m.shared.Close()
		return err
	}

	// SIGHUP is the reload signal. Notify keeps the process from being
	// terminated by it, which is what it would do by default.
	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGHUP)
	m.stop = func() { signal.Stop(reload) }

	go func() {
		for {
			select {
			case <-m.context.Done():
				return

			case <-reload:
				m.logger().Info("caught sighup, reloading")

				if err := m.Reload(ctx); err != nil {
					// Whether to carry on is decided from what is
					// serving, not from the error: Check refuses before
					// retire and leaves the generation in place, while a
					// failure after it leaves current nil.
					if m.Platform() != nil {
						m.logger().Error("reload refused, still serving", "error", err)
						continue
					}

					// Nothing is serving after a failed reload, and a
					// retry would read the same configuration again.
					// Exiting is the honest outcome: it is visible to
					// whatever supervises the process.
					m.logger().Error("reload failed, stopping", "error", err)
					m.Stop()
					return
				}
			}
		}
	}()

	return nil
}

// Reload stops the running platform and starts a new one on the same
// socket. Generations never overlap, so a module registered as a value,
// rather than as a constructor, has to survive a restart.
//
// Check decides whether the swap happens at all. It runs first, while the
// current generation is still serving, so a reload that is refused costs
// nothing: the error comes back and the platform that was answering
// requests goes on answering them.
func (m *Manager) Reload(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Before retire, which is the whole point: past this line the generation
	// that was serving is gone and there is no way back to it.
	if m.Check != nil {
		if err := m.Check(); err != nil {
			return fmt.Errorf("reload refused: %w", err)
		}
	}

	m.retire()

	if err := m.startGeneration(ctx); err != nil {
		return fmt.Errorf("reload: %w", err)
	}

	m.logger().Info("platform reloaded", "url", m.URL())
	return nil
}

// startGeneration builds a platform on the shared socket and starts it.
// The caller holds mu.
func (m *Manager) startGeneration(ctx context.Context) error {
	p := New(m.options)
	p.Logger = m.Logger
	p.listener = m.shared.next()

	// The pidfile records the process, and the manager holds it. Leaving a
	// generation one of its own would have every reload remove the file and
	// write it again, and `-s reload` reads it in exactly that window.
	p.pid = pidfile{}

	if m.Setup != nil {
		if err := m.Setup(p); err != nil {
			p.Stop()
			return fmt.Errorf("manager setup: %w", err)
		}
	}

	if err := p.Start(ctx); err != nil {
		p.Stop()
		return err
	}

	g := &generation{platform: p}
	m.current.Store(g)

	go m.watch(g)

	return nil
}

// retire stops the current generation, if there is one. The caller holds mu.
func (m *Manager) retire() {
	g := m.current.Swap(nil)
	if g == nil {
		return
	}

	g.retired.Store(true)
	g.platform.Stop()
}

// watch takes the manager down with a generation that stopped on its own,
// which is a caught SIGTERM, or the start context being cancelled.
func (m *Manager) watch(g *generation) {
	<-g.platform.Context().Done()

	if g.retired.Load() {
		return
	}

	m.Stop()
}

// Platform returns the generation currently serving, and nil when there is
// none: before Start, after Stop, and after a reload that failed.
func (m *Manager) Platform() *Platform {
	if g := m.current.Load(); g != nil {
		return g.platform
	}
	return nil
}

// URL gives the e2e endpoint URL for requests. A reload does not change it.
func (m *Manager) URL() string {
	return listenerURL(m.shared)
}

// Context returns the cancellation context for the manager. When the context
// finishes, the platform has shut down and the socket is closed.
func (m *Manager) Context() context.Context {
	return m.context
}

// Wait will pause until the manager is stopped. A reload does not end it.
func (m *Manager) Wait() {
	<-m.context.Done()
}

// Stop shuts down the running platform and closes the socket. A stopped
// manager does not start again.
func (m *Manager) Stop() {
	m.once.Do(func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		defer m.cancel()

		m.retire()
		m.stop()

		if err := m.pid.remove(); err != nil {
			m.logger().Error("pidfile", "error", err)
		}

		if m.shared != nil {
			_ = m.shared.Close()
		}
	})
}
