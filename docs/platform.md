# The Platform

## Overview

`platform` is a modular framework for HTTP servers in Go. It provides:

- A global registry for middleware and modules.
- A module lifecycle for graceful startup/shutdown.
- A router (alias to `chi.Router`) for attaching module routes.
- Named database connections with automatic environment scanning.

Each `Platform` instance clones the global registry, enabling isolated test instances and avoiding races or goroutine leaks. The clone calls the constructors registered with `platform.RegisterFunc()`, so each instance holds modules of its own. A value registered with the deprecated `platform.Register()` is shared by every instance in the process.

## Key Concepts

- Module - implements `Name()`, `Start(context.Context)`, `Mount(context.Context, Router)`, `Stop(context.Context)`. Registered as a constructor with `platform.RegisterFunc()`, so each platform builds its own.
- Middleware - type `func(http.Handler) http.Handler`, added via `platform.Use()` or `(*Platform).Use()`.
- Registry - package and instance level container value managing modules and middleware; enables `init` usage via package API.
- Database - named connections, automatically scanned from `PLATFORM_DB_*` environment variables. `"default"` is used if no name is passed.
- Logger - the `Platform.Logger` field, an interface with `Info` and `Error`, receiving the platform's own output.
- Manager - owns the listening socket and the `*Platform` serving on it, replacing the platform on `SIGHUP`. It also owns the pidfile, because a reload does not make a new process.
- Pidfile - the file `Options.PidFile` names, holding the process id something else sends a signal to. Empty writes none.

## Logging

The platform doesn't use the `log` package. It writes through the exported `Platform.Logger` field, declared as:

```go
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
```

`New` sets the field to `slog.Default()`, or to a discarding logger when `Options.Quiet` is set. A `*slog.Logger` satisfies the interface as it is, so a consumer application can hand the platform its own logger:

```go
p := platform.New(platform.NewOptions())
p.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

Assign the field before calling `Start`. The platform reads it once there, and keeps logging through that value for the lifetime of the instance.

Modules reach the same logger from a request or a context:

```go
platform.FromRequest(r).Logger.Info("handled", "path", r.URL.Path)
```

## Lifecycle

1. **Register modules** via `platform.RegisterFunc()` (or `Register` on a `*Platform` instance).
2. **Add middleware** via `platform.Use()` before calling `Start(context.Context)`.
3. **Start the platform** with `Start(context.Context)`; modules are started and then mounted, the socket is bound, and `Options.PidFile` is written when one is named.
4. **Stop** with `Stop()`, which is also what a `SIGINT` or a `SIGTERM` reaches; the server is shut down gracefully with a 5 second timeout, the platform context is cancelled, and the registry then stops every module in parallel.
5. Application exit, reporting any error during shutdown.

## Reload

A `*Platform` is a one-shot value. `Stop` clears its registry and cancels its context, and there is no way back from that, so a reload is a new platform. `Manager` is what outlives the old one and builds the new one:

```go
m := platform.NewManager(platform.NewOptions())
if err := m.Start(ctx); err != nil {
	return err
}
m.Wait()
```

`cmd.Main` runs a manager, so an app built on it reloads with `kill -HUP`. Used directly, `platform.Start` is unchanged, and `SIGHUP` keeps its default disposition, which terminates the process.

The manager holds the listening socket, so a reload keeps the address it was reached on, along with the connections queued on it. Everything above the socket is new: the router, the registry, the server, the telemetry recorder, and the value `Platform()` returns.

Generations do not overlap: the old one is drained and stopped before the new one starts. Requests that arrive during the swap wait in the accept queue of the socket rather than being refused, and requests already in flight are served by the generation that took them.

Modules registered with `platform.RegisterFunc` are constructed per generation, so a reload starts fresh values. A module registered with the deprecated `platform.Register` is one value shared by every generation, and has to tolerate `Start` after `Stop`.

Registrations made against a platform value do not survive a reload, because the value does not. `Manager.Setup` is where they belong:

```go
m.Setup = func(p *platform.Platform) error {
	p.Register(user.NewModule())
	p.Use(middleware.Logger)
	return nil
}
```

A reload that fails leaves nothing serving: the old generation is already gone, and a retry would read the same configuration again. `Reload` returns the error, and the `SIGHUP` handler stops the manager, so the failure is visible to whatever supervises the process.

## Pidfile

`kill -HUP` needs the pid, and `Options.PidFile` is where the process writes it:

```go
options := platform.NewOptions()
options.PidFile = "/run/myapp.pid"
```

The file holds the decimal pid and a newline, created 0644 before the umask, and `platform.ReadPidFile` reads it back. Empty, the default, writes nothing.

It is written once the modules have started and the socket is bound, so the file never names a process that then failed to come up, and removed once the server has drained. A file still present after the process is gone means the process did not stop cleanly.

The manager owns the file, not the generations it runs. A pidfile records a process, and a reload does not make a new one, so a generation neither writes nor removes one: `startGeneration` clears it. Without that, every reload would remove the file and write it again, and the moment a `SIGHUP` sender reads it is exactly the moment it would be missing.

The directory has to exist. It belongs to whatever packages the service, `RuntimeDirectory=` in a systemd unit being the usual case, and creating it here would mean guessing its owner and mode, so a path that names a missing directory fails the start rather than being created. An existing file is overwritten rather than treated as a running instance: a pidfile is a record and not a lock, and refusing to start because of a file a killed process left behind is how a service fails to come back after a power cut.
