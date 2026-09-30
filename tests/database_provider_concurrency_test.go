package platform_test

import (
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/assert"
)

// TestDatabaseProviderSingleton covers the cache holding its lock across the
// open. It used to release between the miss and the build, so every concurrent
// first caller opened a connection of its own: one was cached and the rest were
// handed out unreachable, never closed. On an in-memory DSN they are not even
// the same database.
func TestDatabaseProviderSingleton(t *testing.T) {
	var opens atomic.Int64

	provider := internal.NewDatabaseProvider(func(driver, dsn string) (*sqlx.DB, error) {
		opens.Add(1)
		return sqlx.Open(driver, dsn)
	})
	provider.Register("test", "sqlite://:memory:")

	const callers = 16

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		handles = map[*sqlx.DB]struct{}{}
	)

	wg.Add(callers)
	for range callers {
		go func() {
			defer wg.Done()

			db, err := provider.Open(t.Context(), "test")
			if err != nil {
				return
			}

			mu.Lock()
			handles[db] = struct{}{}
			mu.Unlock()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, len(handles), "every caller must get the same handle")
	assert.Equal(t, int64(1), opens.Load(), "the connection must be opened once")
}

// TestDatabaseProviderCachesUnderResolvedName covers the cache key. A fallback
// list used to be cached under its first name rather than the one that
// resolved, so an unregistered name started answering and the registered one
// built a second handle to the same database.
func TestDatabaseProviderCachesUnderResolvedName(t *testing.T) {
	provider := internal.NewDatabaseProvider(sqlx.Open)
	provider.Register("default", "sqlite://"+filepath.Join(t.TempDir(), "test.db"))

	// "replica" has no credential, so "default" resolves.
	first, err := provider.Open(t.Context(), "replica", "default")
	assert.NoError(t, err)
	assert.NotNil(t, first)

	second, err := provider.Open(t.Context(), "default")
	assert.NoError(t, err)
	assert.Equal(t, first, second, "the resolved name must hit the cache")

	_, err = provider.Open(t.Context(), "replica")
	assert.Error(t, err, "a name with no credential must not be served from the cache")
}

// TestDatabaseProviderEvictsUnpingable covers Connect leaving a handle it could
// not ping in the cache. Open reports no error of its own, because sql.Open
// does not dial, so the dead handle was handed to every later caller with a nil
// error and no way to tell.
func TestDatabaseProviderEvictsUnpingable(t *testing.T) {
	provider := internal.NewDatabaseProvider(sqlx.Open)
	provider.Register("broken", "sqlite://"+filepath.Join(t.TempDir(), "missing", "app.db"))

	first, err := provider.Open(t.Context(), "broken")
	assert.NoError(t, err, "sql.Open does not dial, so it does not fail here")
	assert.NotNil(t, first)

	_, err = provider.Connect(t.Context(), "broken")
	assert.Error(t, err, "a database in a directory that does not exist cannot be pinged")

	second, err := provider.Open(t.Context(), "broken")
	assert.NoError(t, err)
	assert.NotEqual(t, first, second, "the handle whose ping failed must be evicted")

	// Evicting also closes it, so nothing is left holding a pool nobody can
	// reach.
	assert.Error(t, first.PingContext(t.Context()), "the evicted handle is closed")
}
