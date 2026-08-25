package platform

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/titpetric/platform/pkg/require"
)

func newTestSharedListener(tb testing.TB) *sharedListener {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(tb, err)

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
	require.Equal(t, addr, first.Addr().String())

	client, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	conn, err := first.Accept()
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	// Retiring the generation is what http.Server.Shutdown does to the
	// listener it serves. The socket has to survive it.
	require.NoError(t, first.Close())

	_, err = first.Accept()
	require.ErrorIs(t, err, net.ErrClosed)

	// A connection made while no generation is accepting waits in the
	// accept queue of the socket, and the next generation serves it.
	waiting, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = waiting.Close() })

	second := shared.next()
	require.Equal(t, addr, second.Addr().String())

	conn, err = second.Accept()
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	require.NoError(t, shared.Close())
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
	require.NoError(t, first.Close())

	select {
	case err := <-accepted:
		require.ErrorIs(t, err, net.ErrClosed)
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
	require.NoError(t, err)
	require.Equal(t, server, conn)

	t.Run("only one connection can be parked", func(t *testing.T) {
		first, _ := net.Pipe()
		second, _ := net.Pipe()

		shared.handoff(first)
		shared.handoff(second)

		// The slot is taken, so the second connection is closed rather
		// than held for a generation that may never arrive.
		_, err := second.Write([]byte("x"))
		require.ErrorIs(t, err, io.ErrClosedPipe)

		require.NoError(t, shared.Close())

		// Closing the socket releases what was parked on it.
		_, err = first.Write([]byte("x"))
		require.ErrorIs(t, err, io.ErrClosedPipe)
	})
}
