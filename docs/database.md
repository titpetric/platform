# SQL Database Usage

The platform package implements a **named database provider**:

```go
type DatabaseProvider interface {
	Open(ctx context.Context, names ...string) (*sqlx.DB, error)
	Connect(ctx context.Context, names ...string) (*sqlx.DB, error)
}
```

To use, a `platform.Database` value is provided. It's expected for a
module to use a named connection, as an example of a business domain
boundary, and least privilege access.

In practice, a singular modular monolith may share the complete schema
and no named connections need to be used. Passing no name uses the
`"default"` connection, which is all you need on a shared schema:

```go
db, err := platform.Database.Connect(ctx)
```

The platform automatically imports three drivers:

1. `go-sql-driver/mysql` for MySQL, Percona, MariaDB,
2. `github.com/jackc/pgx/v5/stdlib` for PostgreSQL,
3. `modernc.org/sqlite` for sqlite,

The extensions are imported by `pkg/drivers`, integration tests need to
import the package to provide database functionality, or import the
packages (or any other) on their own. The platform is database agnostic.

## Named Connections

The platform scans the runtime environment for `PLATFORM_DB_` prefixed
environment variables. The remainder after the prefix is lowercased and
used for the connection name, so `PLATFORM_DB_USERS` registers `"users"`.

The scan runs from an `init` function, and always registers a `"default"`
connection of `sqlite://:memory:` first. Setting `PLATFORM_DB_DEFAULT`
replaces it.

## Connection strings

```text
sqlite://:memory:
postgres://user:pass@localhost:5432/dbname?sslmode=disable
mysql://user:pass@tcp(localhost:3306)/dbname
```

These are a few connection string examples that can be used to connect
to various databases. The value is constructed as `<driver>://<dsn>`.
Without the `<driver>://` prefix the value is taken as a MySQL DSN.
`postgres` and `postgresql` map onto the `pgx` driver.

The platform fills in driver defaults the DSN does not already set. MySQL
gets `parseTime=true`, `collation=utf8mb4_general_ci` and `loc=Local`.

File-backed SQLite connections default to WAL mode, a 5-second busy timeout,
and a pool of up to 10 open and 2 idle connections. Explicit `_journal_mode`
and `_busy_timeout` DSN options take precedence. In-memory SQLite connections
do not receive these defaults and remain limited to one open and idle connection
so every query uses the same database.

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

The connection does not need to be explicitly closed. A named connection
is reused between modules, the `*sqlx.DB` value returned from repeated
Open or Connect calls will be shared.

The returned database client is safe for concurrent use. Some
restrictions may apply on a per-driver basis.
