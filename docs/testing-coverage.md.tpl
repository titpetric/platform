# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

The storage package has integration tests behind a build tag. To run
integration tests, run `task integration`.

## Packages

{{.Packages}}

## Functions

{{.Functions}}
