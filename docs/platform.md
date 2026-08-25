# The Platform

## Overview

`platform` is a modular framework for HTTP servers in Go. It provides:

- A global registry for middleware and modules.
- A module lifecycle for graceful startup/shutdown.
- A router (alias to `chi.Router`) for attaching module routes.
- Named database connections with automatic environment scanning.

Each `Platform` instance clones the global registry, enabling isolated test instances and avoiding races or goroutine leaks.

## Key Concepts

- Module - implements `Name()`, `Mount(Router)`, `Start(context.Context)`, `Stop()`.
- Middleware - type `func(http.Handler) http.Handler`, added via `platform.Use()` or `(*Platform).Use()`.
- Registry - package and instance level container value managing modules and middleware; enables `init` usage via package API.
- Database - named connections, automatically scanned from `PLATFORM_DB_*` environment variables. `"default"` is used if no name is passed.
- Logger - the `Platform.Logger` field, an interface with `Info` and `Error`, receiving the platform's own output.
- Manager - owns the listening socket and the `*Platform` serving on it, replacing the platform on `SIGHUP`.

## Logging

The platform doesn't use the `log` package. It writes through the exported
`Platform.Logger` field, declared as:

```go
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
```

`New` sets the field to `slog.Default()`, or to a discarding logger when
`Options.Quiet` is set. A `*slog.Logger` satisfies the interface as it is, so
a consumer application can hand the platform its own logger:

```go
p := platform.New(platform.NewOptions())
p.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

Assign the field before calling `Start`. The platform reads it once there,
and keeps logging through that value for the lifetime of the instance.

Modules reach the same logger from a request or a context:

```go
platform.FromRequest(r).Logger.Info("handled", "path", r.URL.Path)
```

## Lifecycle

1. **Register modules** via `platform.Register()` (or on a `*Platform` instance).
2. **Add middleware** via `platform.Use()` before calling `Start(context.Context)`.
3. **Start the platform** with `Start(context.Context)`; modules are started and then mounted.
4. **Stop** with `Stop()`; modules are stopped in parallel, then the server context is cancelled.
5. Application exit, reporting any error during shutdown.

## Reload

A `*Platform` is a one-shot value. `Stop` clears its registry and cancels its
context, and there is no way back from that, so a reload is a new platform.
`Manager` is what outlives the old one and builds the new one:

```go
m := platform.NewManager(platform.NewOptions())
if err := m.Start(ctx); err != nil {
	return err
}
m.Wait()
```

`cmd.Main` runs a manager, so an app built on it reloads with `kill -HUP`.
Used directly, `platform.Start` is unchanged, and `SIGHUP` keeps its default
disposition, which terminates the process.

The manager holds the listening socket, so a reload keeps the address it was
reached on, along with the connections queued on it. Everything above the
socket is new: the router, the registry, the server, the telemetry recorder,
and the value `Platform()` returns.

Generations do not overlap. A module registered from `init` is one value
shared by every generation, cloned out of the global registry, so the old
generation is drained and stopped before the new one starts. For a reload to
mean anything, a module has to tolerate `Start` after `Stop`. Requests that
arrive during the swap wait in the accept queue of the socket rather than
being refused; requests already in flight are served by the generation that
took them.

Registrations made against a platform value do not survive a reload, because
the value does not. `Manager.Setup` is where they belong:

```go
m.Setup = func(p *platform.Platform) error {
	p.Register(user.NewModule())
	p.Use(middleware.Logger)
	return nil
}
```

A reload that fails leaves nothing serving: the old generation is already
gone, and a retry would read the same configuration again. `Reload` returns
the error, and the `SIGHUP` handler stops the manager, so the failure is
visible to whatever supervises the process.
