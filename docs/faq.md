# Frequently asked questions

## How do I register a middleware?

Use `platform.Use` (package) or `(*Platform).Use` (instance). The package function writes to the global registry, which `platform.New` clones once, so call it in `main` ahead of `New`, or in an `init` before `New` or `Start`. The instance method takes effect any time before `(*Platform).Start`. Added too late, middleware is silently never run.

## How do I register a module?

Use `platform.RegisterFunc` (package) or `(*Platform).Register` (instance). The package function must run before `platform.New`, which clones the global registry. The instance method runs any time before the server starts. The package function takes a constructor and calls it once per platform: a reload generation and a parallel test get a module of their own. `platform.Register` takes a value and is deprecated: one value is shared by the platforms in the process.

## How do I access a named database connection?

Use `platform.Database.Connect(ctx, "name")` to Open + Ping a connection, or `Open(ctx, "name")` to skip the ping. Passing no name uses `"default"`. `platform.SetupConnections(os.Environ())` decodes the `PLATFORM_DB_*` variables into the named connections.

The platform imports no sql driver. A blank import in the `main` package registers the one your DSN names:

```go
import (
	_ "github.com/go-sql-driver/mysql"
)
```

Without this, the call fails with `sql: unknown driver "mysql" (forgotten import?)`. See [SQL Database Usage](database.md).

## How do I run the platform in tests?

Create the instance with `platform.New(platform.NewTestOptions())` and call `Register`/`Use` on the instance. Avoid package-level `Register` in tests. The test options bind to `127.0.0.1:0`, silence the platform's own output, and leave telemetry off, so parallel tests do not observe each other.

## How do I implement a module quickly?

Embed `platform.UnimplementedModule` and override only the methods you need.

## How do I start/stop a platform instance?

The package provides `Start(context.Context, *Options)`, a shorthand that allocates a platform and starts it. The options object configures how it starts. `nil` takes the defaults from `NewOptions()`.

The platform is shut down when the context passed to `Start` is cancelled, or when the process is sent a `SIGINT` or a `SIGTERM`.

```go
p, err := platform.Start(ctx, platform.NewOptions())
if err != nil {
	return err
}
```

`Start` returns as soon as the server is accepting. To block until it is done, invoke `Wait()`, like you would with a `sync.WaitGroup`. The function exits when the server has shut down due to cancellation via signal.

```go
p := platform.New(platform.NewOptions())

if err := p.Start(ctx); err != nil {
	return err
}

p.Wait()
return nil
```

The alternative to `p.Wait()` is to use `p.Stop()` explicitly when you want to shut down the platform.

## How do I attach routes in a module?

Implement `Mount(ctx context.Context, r Router) error` to register GET/POST handlers via `r.Get()`/`r.Post()` and other options. Functions exist to add grouping to your endpoints, like `r.Route(prefix, func(Router))`. This gives you simple ways to use middleware in your routes.

## How do I handle graceful shutdown?

Graceful shutdown is implemented by the platform. In the modules you need to implement `Start` and `Stop` functions, which should create and cancel any goroutines needed by your module.

The platform will shut itself down if a `SIGINT` or a `SIGTERM` is caught. `SIGKILL` cannot be caught by any process, so nothing runs. For testing, the passed context is expected to be a `t.Context()` (for `testing.TB`).

Start has a platform instance attached to the context, and can use `platform.FromContext` to get the instance, and the `Find` function on the instance to get a reference to any side loaded module.

## How do I reload the app without restarting the process?

Run a `platform.Manager` instead of a bare `*Platform`, and send the process a `SIGHUP`. `cmd.Main` already does, and the app built on the manager reloads out of the box:

```go
m := platform.NewManager(platform.NewOptions())
if err := m.Start(ctx); err != nil {
	return err
}

m.Wait()
```

The manager holds the listening socket and replaces the platform under it, and the address survives the reload while the router, the registry and the modules are new. Modules registered with `RegisterFunc` are constructed per generation. A module registered as a value with the deprecated `Register` is shared across generations and must tolerate `Start` after `Stop`. Registrations that are not in the global registry belong in `Manager.Setup`, which runs against every generation. See the reload section in [The Platform](platform.md).

## How do I stop a bad configuration from taking the server down on reload?

Set `Manager.Check`. It runs before the reload retires the platform that is serving, and a non-nil error abandons the reload with that platform left answering requests:

```go
m.Check = func() error {
	_, err := config.Load(filename)
	return err
}
```

Without it, a `SIGHUP` carrying a configuration that fails to load stops the old generation first and then fails to build the new one, which leaves nothing serving and stops the process.

Assign it ahead of `Start`, as with `Setup`: the signal handler reads it on a separate goroutine.

## How does something else find the process to signal it?

Set `Options.PidFile`, or `PLATFORM_PIDFILE`, and the process records its id there on start and removes it on a clean stop:

```go
options := platform.NewOptions()
options.PidFile = "/run/myapp.pid"
```

The file holds the decimal pid and a newline, which is the form `kill -HUP $(cat /run/myapp.pid)` and all service managers expect. The platform exports no reader for it: whatever sends the signal already has the path from the configuration.

Empty, the default, writes no file at all. The directory must exist: it belongs to whatever packages the service, and creating it on this side would mean guessing its owner and mode. An existing file is overwritten rather than treated as a running instance, because a pidfile is a record and not a lock.

A `Manager` writes the file once for the process, and the generations it runs do not, and the reload leaves it alone. A file still present after the process is gone means the process did not stop cleanly.

`SIGHUP` reloads only under a `Manager`. Sent to a process running a bare `platform.Start`, it keeps its default disposition and terminates it. `SIGTERM` stops either.
