# SQL Database Usage

The platform package implements a **named database provider**:

```go
type DatabaseProvider interface {
	Open(ctx context.Context, names ...string) (*sqlx.DB, error)
	Connect(ctx context.Context, names ...string) (*sqlx.DB, error)
}
```

To use, a `platform.Database` value is provided. It's expected for a module to use a named connection, as an example of a business domain boundary, and least privilege access.

In practice, a singular modular monolith may share the complete schema and no named connections need to be used. Passing no name uses the `"default"` connection, which is all you need on a shared schema:

```go
db, err := platform.Database.Connect(ctx)
```

`Open` returns a handle without contacting the server, as `sql.Open` does. `Connect` is `Open` plus a `PingContext`, so it is the call that fails on a bad DSN, an unreachable server or a driver that was never registered. Both return the same cached handle for a given name.

The platform imports no driver. It speaks to `database/sql` through `sqlx`, and the driver is the choice of whoever builds the binary. Register one with a blank import in your `main` package:

```go
import (
	_ "github.com/go-sql-driver/mysql" // MySQL, Percona, MariaDB
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL
	_ "modernc.org/sqlite"             // sqlite
)
```

Any driver that registers itself with `database/sql` works, under whatever name it registers. The DSN prefix selects it; see [Connection strings](#connection-strings).

Pinning a driver in the platform would push that pin onto every consumer and duplicate the dependency tree of anyone already on a newer version. The tests in `tests/` are a separate Go module for the same reason: they register all three drivers in `tests/main_test.go`, and none of that reaches the platform `go.mod`.

## Named Connections

The platform scans the runtime environment for `PLATFORM_DB_` prefixed environment variables. The variable name after the prefix is lowercased and used for the connection name; the value is the connection string, taken verbatim. `PLATFORM_DB_USERS` registers `"users"`.

The scan runs from an `init` function over `os.Environ()`. The `"default"` connection is seeded with `sqlite://:memory:` before the environment is read, so `PLATFORM_DB_DEFAULT` overrides it and a process with nothing set still has a usable default. Registration order across names is unspecified.

`platform.SetupConnections([]string)` runs the same scan against an environment you supply. `platform.Database` exposes only `Open` and `Connect`, so that call and the `PLATFORM_DB_*` variables are the whole configuration surface; a consumer wanting different behaviour assigns its own implementation to `platform.Database`.

## Connection strings

```text
sqlite://:memory:
postgres://user:pass@localhost:5432/dbname?sslmode=disable
mysql://user:pass@tcp(localhost:3306)/dbname
```

These are a few connection string examples that can be used to connect to various databases. The value is constructed as `<driver>://<dsn>`. Without the `<driver>://` prefix the value is taken as a MySQL DSN. `postgres` and `postgresql` map onto the `pgx` driver, and the scheme is put back on the DSN handed to it, so both the URL form and a libpq keyword string work unmodified.

The platform fills in driver defaults the DSN does not already set. MySQL gets `parseTime=true`, `collation=utf8mb4_general_ci` and `loc=Local`, and a pool of up to 10 open and 10 idle connections.

The DSN defaults key on the exact driver names `sqlite` and `mysql`. Any other driver, `pgx` included, gets none of them, and a pool of up to 10 open and 2 idle connections. Set the pool yourself on the returned `*sqlx.DB` if that is not what you want. A driver registered under a different name, `sqlite3` rather than `sqlite`, opens fine and gets the same fallback, which for an in-memory DSN means the one-connection cap is not applied and each pooled connection sees a database of its own.

File-backed SQLite connections default to WAL mode, a 5-second busy timeout, and a pool of up to 10 open and 2 idle connections. Explicit `_journal_mode` and `_busy_timeout` DSN options take precedence. In-memory SQLite connections do not receive these defaults and remain limited to one open and idle connection so every query uses the same database. A DSN counts as in-memory when its path is `:memory:` or `file::memory:`, or when its query sets `mode=memory`.

## Using Connections in Modules

```go
func (m *Module) Start(ctx context.Context) error {
	db, err := platform.Database.Connect(ctx) // Open + Ping
	if err != nil {
		return err
	}
	m.storage = NewStorage(db)
	return nil
}
```

The connection does not need to be explicitly closed. A named connection is reused between modules: repeated `Open` or `Connect` calls with the same name return the same `*sqlx.DB`. Passing several names is a fallback list, and the first name with a registered credential is the one opened. The handle is cached under that name, so two callers with different fallback lists resolving to the same name share one handle.

A `Connect` whose ping fails reports the error and leaves the handle cached. `*sqlx.DB` is a pool and reconnects on its own, so a name that was unreachable once serves the same handle when the database comes back.

`platform.Transaction(ctx, db, fn)` runs `fn` in a transaction, committing when it returns nil and rolling back on an error or a panic.

The returned database client is safe for concurrent use. Some restrictions may apply on a per-driver basis.
