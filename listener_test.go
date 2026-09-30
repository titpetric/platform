package platform

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/titpetric/platform/internal/assert"
)

func newTestSharedListener(tb testing.TB) *sharedListener {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(tb, err)

	shared := newSharedListener(listener)
	tb.Cleanup(func() { _ = shared.Close() })

	return shared
}

// TestSharedListener covers the socket outliving the generation serving on
// it, which is what makes a reload keep its address.
func TestSharedListener(t *testing.T) {
	shared := newTestSharedListener(t)
	addr := shared.Addr().String()

	first := shared.next()
	assert.Equal(t, addr, first.Addr().String())

	client, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	conn, err := first.Accept()
	assert.NoError(t, err)
	assert.NoError(t, conn.Close())

	// Retiring the generation is what http.Server.Shutdown does to the
	// listener it serves. The socket has to survive it.
	assert.NoError(t, first.Close())

	_, err = first.Accept()
	assert.ErrorIs(t, err, net.ErrClosed)

	// A connection made while no generation is accepting waits in the
	// accept queue of the socket, and the next generation serves it.
	waiting, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	t.Cleanup(func() { _ = waiting.Close() })

	second := shared.next()
	assert.Equal(t, addr, second.Addr().String())

	conn, err = second.Accept()
	assert.NoError(t, err)
	assert.NoError(t, conn.Close())

	assert.NoError(t, shared.Close())
}

// TestSharedListener_unblocks covers the accept loop of a generation that
// is retired while it is blocked on the socket it does not own.
func TestSharedListener_unblocks(t *testing.T) {
	shared := newTestSharedListener(t)

	first := shared.next()

	accepted := make(chan error, 1)
	go func() {
		_, err := first.Accept()
		accepted <- err
	}()

	// Give the accept loop a chance to block on the socket.
	time.Sleep(10 * time.Millisecond)
	assert.NoError(t, first.Close())

	select {
	case err := <-accepted:
		assert.ErrorIs(t, err, net.ErrClosed)
	case <-time.After(5 * time.Second):
		t.Fatal("accept did not return after the generation was retired")
	}
}

// TestSharedListener_handoff covers the connection a retired generation
// accepted but never served. It belongs to the next generation.
func TestSharedListener_handoff(t *testing.T) {
	shared := newTestSharedListener(t)

	server, client := net.Pipe()
	t.Cleanup(func() { _ = client.Close() })

	shared.handoff(server)

	conn, err := shared.next().Accept()
	assert.NoError(t, err)
	assert.Equal(t, server, conn)

	t.Run("only one connection can be parked", func(t *testing.T) {
		first, _ := net.Pipe()
		second, _ := net.Pipe()

		shared.handoff(first)
		shared.handoff(second)

		// The slot is taken, so the second connection is closed rather
		// than held for a generation that may never arrive.
		_, err := second.Write([]byte("x"))
		assert.ErrorIs(t, err, io.ErrClosedPipe)

		assert.NoError(t, shared.Close())

		// Closing the socket releases what was parked on it.
		_, err = first.Write([]byte("x"))
		assert.ErrorIs(t, err, io.ErrClosedPipe)
	})
}

// TestListenerURLNoSocket covers the nil dereference listenerURL used to do. A
// platform or manager that has not started, or whose start failed, has no
// socket, and logging the failure with the URL in hand panicked.
func TestListenerURLNoSocket(t *testing.T) {
	assert.Equal(t, "", listenerURL(nil))

	// Platform holds a nil net.Listener interface before Start. Manager holds
	// a nil *sharedListener inside a non-nil one, which is why Manager.URL has
	// a check of its own.
	assert.Equal(t, "", New(NewTestOptions()).URL())
	assert.Equal(t, "", NewManager(NewTestOptions()).URL())
}
