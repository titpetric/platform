# Structural diagram

## Modules

The repository holds two Go modules.

`github.com/titpetric/platform` is the module consumers import. It depends on chi, sqlx and oida, and on no sql driver. Its tests are the ones that read unexported symbols, so they have to sit in the package they test and they add no dependency of their own.

`github.com/titpetric/platform/tests` holds the tests that reach the platform through its exported API. It replaces the platform with the checkout above it, and it is where the mysql, pgx and sqlite drivers are pinned. `tests/main_test.go` registers all three with `database/sql` in package init, so every test in the suite can open any of the three DSN forms without importing a driver.

The suite is an external test package, not a black box: it imports `internal/assert` for its assertions, and it covers `internal`, `internal/httpcontext` and `internal/pidfile` through their exported API. `internal/httpcontext` has no test left in the root module at all.

The split keeps driver pins out of the consumer dependency tree. Run both with `atkins`; the root module alone is `go test ./...`, the suite is `go test ./...` from `tests/`.

## Diagrams

Both diagrams are generated from source by `atkins diagrams`, which the `cover` job runs. Do not edit the `.puml` or `.svg` files by hand.

Import diagram:

![](imports.svg)

Class diagram:

![](structure.svg)
