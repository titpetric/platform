# The Platform

## Overview

`platform` is a modular framework for HTTP servers in Go. It provides:

- A **global registry** for middleware and modules.
- A **module lifecycle** for graceful startup/shutdown.
- A **router** (alias to `chi.Router`) for attaching module routes.
- **Named database connections** with automatic environment scanning.

Each `Platform` instance clones the global registry, enabling isolated test instances and avoiding races or goroutine leaks.

## Key Concepts

- **Module** — implements `Name()`, `Mount(Router)`, `Start()`, `Stop()`.
- **Middleware** — type `func(http.Handler) http.Handler`, added via `platform.Use()` or `(*Platform).Use()`.
- **Registry** — package and instance level container value managing modules and middleware; enables `init` usage via package API.
- **Database** — named connections, automatically scanned from `PLATFORM_DB_*` environment variables. `"default"` is used if no name is passed.

## Lifecycle

1. **Register modules** via `platform.Register()` (or on a `*Platform` instance).
2. **Add middleware** via `platform.Use()` before calling `Start()`.
3. **Start the platform** with `Start(context.Context)`; modules are started and then mounted.
4. **Stop** with `Stop()`; modules are stopped in parallel, then the server context is cancelled.
5. Application exit, reporting any error during shutdown.
