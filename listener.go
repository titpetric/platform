package platform

import (
	"net"
	"sync"
	"time"
)

// sharedListener is a socket that outlives each platform serving it.
type sharedListener struct {
	net.Listener

	// pending holds a connection accepted by a generation that lost the
	// socket before it could serve it. One slot is enough: an http.Server
	// has one accept loop, so one Accept can be in flight.
	pending chan net.Conn
}

// deadliner unblocks an accept loop without closing the socket.
type deadliner interface {
	SetDeadline(t time.Time) error
}

func newSharedListener(l net.Listener) *sharedListener {
	return &sharedListener{
		Listener: l,
		pending:  make(chan net.Conn, 1),
	}
}

// next hands the socket to a new generation, clearing the old deadline.
func (s *sharedListener) next() *generationListener {
	if d, ok := s.Listener.(deadliner); ok {
		_ = d.SetDeadline(time.Time{})
	}

	return &generationListener{
		shared: s,
		closed: make(chan struct{}),
	}
}

// handoff parks a connection for the next generation, closing it when the slot is full.
func (s *sharedListener) handoff(conn net.Conn) {
	select {
	case s.pending <- conn:
	default:
		_ = conn.Close()
	}
}

// Close closes the socket, and anything parked there.
func (s *sharedListener) Close() error {
	select {
	case conn := <-s.pending:
		_ = conn.Close()
	default:
	}

	return s.Listener.Close()
}

// generationListener is one platform generation's view of the shared socket.
type generationListener struct {
	shared *sharedListener

	once   sync.Once
	closed chan struct{}
}

var _ net.Listener = (*generationListener)(nil)

// Accept returns the next connection, or net.ErrClosed once retired.
func (l *generationListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.shared.pending:
		return conn, nil
	default:
	}

	conn, err := l.shared.Accept()

	select {
	case <-l.closed:
		// The socket moved on while this accept was blocked there, so
		// the connection is the next generation's to serve.
		if conn != nil {
			l.shared.handoff(conn)
		}
		return nil, net.ErrClosed
	default:
	}

	return conn, err
}

// Close retires the generation without closing the shared socket.
func (l *generationListener) Close() error {
	l.once.Do(func() {
		close(l.closed)

		if d, ok := l.shared.Listener.(deadliner); ok {
			_ = d.SetDeadline(time.Now())
		}
	})

	return nil
}

// Addr returns the address of the shared socket.
func (l *generationListener) Addr() net.Addr {
	return l.shared.Addr()
}

// listenerURL gives the e2e URL of a listener, or "" when it has no address.
func listenerURL(l net.Listener) string {
	if l == nil {
		return ""
	}

	addr := l.Addr()
	if addr == nil {
		return ""
	}

	_, port, _ := net.SplitHostPort(addr.String())
	return "http://127.0.0.1:" + port
}
