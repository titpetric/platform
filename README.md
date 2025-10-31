# Platform Documentation

## Motivation

The `platform` package is an extensible, modular system for building
HTTP servers and sidecar services in Go.

It provides a global registry for modules and middleware, a lifecycle
for graceful shutdown, and named database connections, allowing you to
structure services as composable, testable modules.

## Development docs

- [The Platform](platform.md) — key concepts, structure, and lifecycle overview.
- [Creating Modules](modules.md) — module API, lifecycle, and using `UnimplementedModule`.
- [Common Patterns](patterns.md) — routing, GET/POST, background jobs and validation patterns.
- [SQL Database Usage](database.md) — named connections, DSN examples, and `Connect()` vs `Open()`.
- [FAQ](faq.md) — short practical answers to common questions.
