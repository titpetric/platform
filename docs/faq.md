# Frequently asked questions

## How do I register a middleware?

Use `platform.Use` (package) or `(*Platform).Use` (instance). Add before calling `Start(context.Context)`.

## How do I register a module?

Use `platform.RegisterFunc` (package) or `(*Platform).Register` (instance)
before starting the server. The package function takes a constructor, and
calls it once per platform, so a reload generation and a parallel test get
a module of their own. `platform.Register` takes a value and is deprecated:
one value is shared by every platform in the process.

## How do I access a named database connection?

Use `platform.Database.Connect("name")` to Open + Ping a connection, or `Open("name")` for lazy access.

## How do I run the platform in tests?

Create a `*Platform` instance and call `Register`/`Use` on it. Avoid package-level `Register` in tests.

## How do I implement a module quickly?

Embed `platform.UnimplementedModule` and override only the methods you need.

## How do I start/stop a platform instance?

The package provides a `Start(context.Context)` function, wrapping a
`StartPlatform(context.Context, *Options)`. Using `StartPlatform` allows
to configure how the platform starts by providing an options object.

The platform is shut down when the context passed to `StartPlatform` is
cancelled or when a SIGTERM signal is intercepted in the system.

```go
if err := platform.Start(ctx); err != nil {
	return log.Fatalf("exit error: %v", err)
}
```

For the `*Platform` instance, you need to invoke the `Wait()` function,
like you would with a `sync.WaitGroup`. The function exits when the
server has shut down due to cancellation via signal.

```go
p, err := platform.New()
if err != nil {
	return err
}

if err := p.Start(ctx); err != nil {
	return err
}

p.Wait()
return nil
```

The alternative to `p.Wait()` is to use `p.Stop()` explicitly when you
want to shut down the platform.

## How do I attach routes in a module?

Implement your `Mount(r Router)` to register GET/POST handlers via
`r.Get()`/`r.Post()` and other options. Functions exist to add grouping
to your endpoints, like `r.Route(prefix, func(Router))`. This gives you
simple ways to use middleware in your routes.

## How do I handle graceful shutdown?

Graceful shutdown is implemented by the platform. In your modules you
need to implement `Start` and `Stop` functions, which should create and
cancel any goroutines needed by your module.

The platform will shut itself down if a `SIGTERM` is caught. For testing,
the passed context is expected to be a `t.Context()` (for `testing.TB`).

Start has a platform instance attached to the context, and can use
`platform.FromContext` to get the instance, and the `Find` function on
the instance to get a reference to any side loaded module.

## How do I reload the app without restarting the process?

Run a `platform.Manager` instead of a bare `*Platform`, and send the
process a `SIGHUP`. `cmd.Main` already does, so an app built on it
reloads out of the box:

```go
m := platform.NewManager(platform.NewOptions())
if err := m.Start(ctx); err != nil {
	return err
}

m.Wait()
```

The manager holds the listening socket and replaces the platform under
it, so the address survives the reload while the router, the registry
and the modules are new. Modules therefore have to tolerate `Start`
after `Stop`. Registrations that are not in the global registry belong
in `Manager.Setup`, which runs against every generation. See the reload
section in [The Platform](platform.md).
