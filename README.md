# platform - A modular system for building Go applications

## Motivation

The `platform` package is an extensible, modular system for building HTTP servers and sidecar services in Go.

It provides a global registry for modules and middleware, a lifecycle for graceful shutdown, and named database connections, allowing you to structure services as composable, testable modules. It imports no sql driver: the binary registers the one its DSN names with a blank import.

Running a `platform.Manager` rather than a bare `*Platform` adds a `SIGHUP` reload that replaces the platform without dropping the socket. `cmd.Main` already does. A bare `platform.Start` leaves `SIGHUP` at its default disposition, which terminates the process.

Application examples, with database use:

- A monolithic app with modules: [titpetric/platform-app](https://github.com/titpetric/platform-app).
- Extended `titpetric/platform-app` for a Mailing list manager app: [titpetric/platform-maillist](https://github.com/titpetric/platform-maillist).

Status: the app and maillist repositories still need implementation surface.

## Coverage

| Status | Package                                 | Coverage | Cognitive | Lines |
|--------|-----------------------------------------|----------|-----------|-------|
| ✅     | titpetric/platform                      | 93.65%   | 146       | 1057  |
| ✅     | titpetric/platform/cmd                  | 46.67%   | 2         | 20    |
| ✅     | titpetric/platform/cmd/platform         | 0.00%    | 0         | 3     |
| ✅     | titpetric/platform/internal             | 88.51%   | 46        | 246   |
| ✅     | titpetric/platform/internal/assert      | 98.57%   | 67        | 270   |
| ✅     | titpetric/platform/internal/httpcontext | 100.00%  | 1         | 18    |
| ✅     | titpetric/platform/internal/pidfile     | 95.83%   | 11        | 44    |

For more detail, see: [Testing Coverage](./docs/testing-coverage.md).

## Development docs

- [The Platform](./docs/platform.md) - key concepts, logging, lifecycle, the SIGHUP reload and the pidfile.
- [API documentation](./docs/api.md) - generated api documentation for the platform package.
- [Creating Modules](./docs/modules.md) - module API, lifecycle, and using `UnimplementedModule`.
- [Common Patterns](./docs/patterns.md) - routing, GET/POST, background jobs, middleware and validation patterns.
- [SQL Database Usage](./docs/database.md) - named connections, DSN examples, and bringing your own sql driver.
- [Telemetry](./docs/telemetry.md) - recording traces and spans with oida, and the `/debug/oida` dashboard.
- [Structural diagram](./docs/structure.md) - the two Go modules, and generated package import and class diagrams.
- [FAQ](./docs/faq.md) - short practical answers to common questions.
