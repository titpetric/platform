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
| 0  | ✅      | github.com/titpetric/platform                     | 84.61%   | 34        | 311   |
| 1  | ✅      | github.com/titpetric/platform/cmd/platform        | 41.65%   | 1         | 16    |
| 2  | ✅      | github.com/titpetric/platform/internal            | 80.18%   | 26        | 184   |
| 3  | ✅      | github.com/titpetric/platform/internal/reflect    | 100.00%  | 7         | 31    |
| 4  | ✅      | github.com/titpetric/platform/internal/require    | 0.00%    | 0         | 0     |
| 5  | ✅      | github.com/titpetric/platform/module/autoload     | 100.00%  | 0         | 4     |
| 6  | ✅      | github.com/titpetric/platform/module/expvar       | 95.23%   | 1         | 19    |
| 7  | ✅      | github.com/titpetric/platform/module/theme        | 100.00%  | 0         | 11    |
| 8  | ❌      | github.com/titpetric/platform/module/user         | 60.41%   | 13        | 103   |
| 9  | ✅      | github.com/titpetric/platform/module/user/model   | 8.11%    | 1         | 82    |
| 10 | ❌      | github.com/titpetric/platform/module/user/service | 28.25%   | 27        | 285   |
| 11 | ❌      | github.com/titpetric/platform/module/user/storage | 20.00%   | 24        | 200   |
| 12 | ❌      | github.com/titpetric/platform/telemetry           | 70.90%   | 8         | 130   |

For more detail, see: [Testing Coverage](./docs/testing-coverage.md).

## Development docs

- [The Platform](./docs/platform.md) — key concepts, structure, and lifecycle overview.
- [Creating Modules](./docs/modules.md) — module API, lifecycle, and using `UnimplementedModule`.
- [Common Patterns](./docs/patterns.md) — routing, GET/POST, background jobs and validation patterns.
- [SQL Database Usage](./docs/database.md) — named connections, DSN examples, and `Connect()` vs `Open()`.
- [Telemetry](./docs/telemetry.md) - setting up and using OpenTelemetry.
- [FAQ](./docs/faq.md) — short practical answers to common questions.
