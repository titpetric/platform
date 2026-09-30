package platform_test

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/assert"
)

// TestDatabaseProviderSingleton covers concurrent first callers getting one
// handle from one open, which an in-memory DSN depends on: two handles are two
// databases.
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

// TestDatabaseProviderCachesUnderResolvedName covers the cache key of a
// fallback list: the name whose credential resolved, not the first requested.
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

// TestDatabaseProviderUnpingableRecovers covers a Connect whose ping fails. The
// handle is a pool and reconnects on its own, so it stays cached and serves the
// callers already holding it once the database is reachable.
func TestDatabaseProviderUnpingableRecovers(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")

	provider := internal.NewDatabaseProvider(sqlx.Open)
	provider.Register("broken", "sqlite://"+filepath.Join(missing, "app.db"))

	first, err := provider.Open(t.Context(), "broken")
	assert.NoError(t, err, "sql.Open does not dial, so it does not fail here")
	assert.NotNil(t, first)

	_, err = provider.Connect(t.Context(), "broken")
	assert.Error(t, err, "a database in a directory that does not exist cannot be pinged")

	assert.NoError(t, os.MkdirAll(missing, 0o755))

	second, err := provider.Connect(t.Context(), "broken")
	assert.NoError(t, err, "the same handle pings once the database is reachable")
	assert.Equal(t, first, second, "a failed ping must not replace the cached handle")

	// The caller that took the handle before the failure still owns a usable
	// one, which closing it on eviction would have taken away.
	assert.NoError(t, first.PingContext(t.Context()))
}
