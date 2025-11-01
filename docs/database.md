# SQL Database Usage

The platform package implements **named database provider**:

```go
type DatabaseProvider interface {
	Add(name, dsn string)
	Open(...string) (*sqlx.DB, error)
	Connect(...string) (*sqlx.DB, error)
}
```

To use, a `platform.Database` value is provided. It's expected for a
module to use a named connection, as an example of a business domain
boundary, and least privilege access.

In practice, a singular modular monolith may share the complete schema
and no named connections need to be used. In your modules `Start`
function you only need this if you're working on shared schema:

```go
db, err := platform.Database.Get()
```

## Named Connections

The platform scans the runtime environment for `PLATFORM_DB_` prefixed
environment variables. The remainder after the prefix is used for the
connection name.

This is done from `init`, and automatically invokes `plaform.Database.Add`.

To register additional connection strings at runtime, you should invoke
the `Add` function. This is the intended pattern if you want to provide
configuration from a config file or other location.

## Connection strings

```text
sqlite://:memory:
postgres://user:pass@localhost:5432/dbname?sslmode=disable
mysql://user:pass@tcp(localhost:3306)/dbname
```

These are a few connection string examples that can be used to connect
to various databases. The value is constructed as `<driver>://<dsn>`.
The platform may add some flags automatically if required by the driver,
like `parseTime=true`.

## Using Connections in Modules

```go
func (m *Module) Start() error {
	db, err := platform.Database.Connect() // Open + Ping
	if err != nil {
		return err
	}
	m.storage = NewStorage(db)
	return nil
}
```

The connection does not need to be explicitly closed. A named connection
is reused between modules, the `*sql.DB` value returned from repeated
Open or Connect calls will be shared.

The returned database client is safe for concurrent use. Some
restrictions may apply on a per-driver basis.
