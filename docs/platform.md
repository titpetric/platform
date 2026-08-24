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
