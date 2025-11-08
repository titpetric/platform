# Platform Documentation

## Motivation

The `platform` package is an extensible, modular system for building
HTTP servers and sidecar services in Go.

It provides a global registry for modules and middleware, a lifecycle
for graceful shutdown, and named database connections, allowing you to
structure services as composable, testable modules.

Application examples, with database use:

- A monolithic app with modules: [titpetric/platform-app](https://github.com/titpetric/platform-app).
- Extended `titpetric/platform-app` for a Mailing list manager app: [titpetric/platform-maillist](https://github.com/titpetric/platform-maillist).
- Miscelaneous module examples: [titpetric/platform-example](https://github.com/titpetric/platform-example).

Status: the app and maillist packages still need implementation surface.

## Coverage

| Status | Package                                     | Coverage | Cognitive | Lines |
| ------ | ------------------------------------------- | -------- | --------- | ----- |
| ✅      | titpetric/platform               | 82.20%   | 42        | 370   |
| ✅      | titpetric/platform/cmd           | 87.50%   | 2         | 23    |
| ✅      | titpetric/platform/cmd/platform  | 0.00%    | 0         | 3     |
| ✅      | titpetric/platform/internal      | 94.87%   | 20        | 155   |
| ✅      | titpetric/platform/pkg/assert    | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform/pkg/drivers   | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform/pkg/reflect   | 100.00%  | 7         | 31    |
| ✅      | titpetric/platform/pkg/require   | 0.00%    | 0         | 0     |
| ❌      | titpetric/platform/pkg/telemetry | 57.44%   | 8         | 130   |
| ✅      | titpetric/platform/pkg/ulid      | 100.00%  | 0         | 20    |

For more detail, see: [Testing Coverage](./docs/testing-coverage.md).

## Development docs

- [The Platform](./docs/platform.md) — key concepts, structure, and lifecycle overview.
- [API documentation](./docs/api.md) - the api documentation for the platform package.
- [Creating Modules](./docs/modules.md) — module API, lifecycle, and using `UnimplementedModule`.
- [Common Patterns](./docs/patterns.md) — routing, GET/POST, background jobs and validation patterns.
- [SQL Database Usage](./docs/database.md) — named connections, DSN examples, and `Connect()` vs `Open()`.
- [Telemetry](./docs/telemetry.md) - setting up and using OpenTelemetry.
- [FAQ](./docs/faq.md) — short practical answers to common questions.
