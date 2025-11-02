# Platform Documentation

## Motivation

The `platform` package is an extensible, modular system for building
HTTP servers and sidecar services in Go.

It provides a global registry for modules and middleware, a lifecycle
for graceful shutdown, and named database connections, allowing you to
structure services as composable, testable modules.

## Coverage

| #  | Status | Package                                           | Coverage | Cognitive | Lines |
| -- | ------ | ------------------------------------------------- | -------- | --------- | ----- |
| 0  | ✅      | titpetric/platform                     | 84.61%   | 34        | 311   |
| 1  | ✅      | titpetric/platform/cmd/platform        | 41.65%   | 1         | 16    |
| 2  | ✅      | titpetric/platform/internal            | 80.66%   | 22        | 192   |
| 3  | ✅      | titpetric/platform/internal/reflect    | 100.00%  | 7         | 31    |
| 4  | ✅      | titpetric/platform/internal/require    | 0.00%    | 0         | 0     |
| 5  | ✅      | titpetric/platform/module/autoload     | 100.00%  | 0         | 4     |
| 6  | ✅      | titpetric/platform/module/expvar       | 95.23%   | 1         | 19    |
| 7  | ✅      | titpetric/platform/module/theme        | 100.00%  | 0         | 11    |
| 8  | ❌      | titpetric/platform/module/user         | 60.41%   | 13        | 103   |
| 9  | ✅      | titpetric/platform/module/user/model   | 8.11%    | 1         | 82    |
| 10 | ❌      | titpetric/platform/module/user/service | 28.25%   | 27        | 285   |
| 11 | ❌      | titpetric/platform/module/user/storage | 20.00%   | 24        | 200   |
| 12 | ❌      | titpetric/platform/telemetry           | 70.90%   | 8         | 130   |

For more detail, see: [Testing Coverage](./docs/testing-coverage.md).

## Development docs

- [The Platform](./docs/platform.md) — key concepts, structure, and lifecycle overview.
- [Creating Modules](./docs/modules.md) — module API, lifecycle, and using `UnimplementedModule`.
- [Common Patterns](./docs/patterns.md) — routing, GET/POST, background jobs and validation patterns.
- [SQL Database Usage](./docs/database.md) — named connections, DSN examples, and `Connect()` vs `Open()`.
- [Telemetry](./docs/telemetry.md) - setting up and using OpenTelemetry.
- [FAQ](./docs/faq.md) — short practical answers to common questions.
