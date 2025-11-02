# Platform Documentation

## Motivation

The `platform` package is an extensible, modular system for building
HTTP servers and sidecar services in Go.

It provides a global registry for modules and middleware, a lifecycle
for graceful shutdown, and named database connections, allowing you to
structure services as composable, testable modules.

## Coverage

| Status | Package                                           | Coverage | Cognitive | Lines |
| ------ | ------------------------------------------------- | -------- | --------- | ----- |
| ✅      | titpetric/platform                     | 84.61%   | 34        | 311   |
| ✅      | titpetric/platform/cmd/platform        | 41.65%   | 1         | 16    |
| ✅      | titpetric/platform/drivers             | 0.00%    | 0         | 0     |
| ❌      | titpetric/platform/internal            | 79.59%   | 24        | 203   |
| ✅      | titpetric/platform/internal/reflect    | 100.00%  | 7         | 31    |
| ✅      | titpetric/platform/internal/require    | 0.00%    | 0         | 0     |
| ✅      | titpetric/platform/module/autoload     | 100.00%  | 0         | 4     |
| ✅      | titpetric/platform/module/expvar       | 95.23%   | 1         | 19    |
| ✅      | titpetric/platform/module/theme        | 100.00%  | 0         | 11    |
| ❌      | titpetric/platform/module/user         | 54.76%   | 13        | 98    |
| ✅      | titpetric/platform/module/user/model   | 10.53%   | 3         | 88    |
| ❌      | titpetric/platform/module/user/service | 28.25%   | 27        | 285   |
| ❌      | titpetric/platform/module/user/storage | 44.01%   | 24        | 211   |
| ❌      | titpetric/platform/telemetry           | 78.59%   | 8         | 130   |

For more detail, see: [Testing Coverage](./docs/testing-coverage.md).

## Development docs

- [The Platform](./docs/platform.md) — key concepts, structure, and lifecycle overview.
- [Creating Modules](./docs/modules.md) — module API, lifecycle, and using `UnimplementedModule`.
- [Common Patterns](./docs/patterns.md) — routing, GET/POST, background jobs and validation patterns.
- [SQL Database Usage](./docs/database.md) — named connections, DSN examples, and `Connect()` vs `Open()`.
- [Telemetry](./docs/telemetry.md) - setting up and using OpenTelemetry.
- [FAQ](./docs/faq.md) — short practical answers to common questions.
