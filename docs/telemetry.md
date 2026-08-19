# Telemetry

The platform records telemetry with [oida](https://github.com/titpetric/oida): traces and spans held in a ring buffer inside the process, with a dashboard served at `/debug/oida`. There is no collector, no exporter and no second service to run.

`platform.New` registers the recorder and the dashboard by default, so a service gets both without wiring anything.

## Instrumenting your code

Import `github.com/titpetric/oida` and start a span. Everything is nil-safe: the same code runs unchanged in a process with no recorder, and in a request that was not sampled.

```go
func (s *UserStorage) GetUsers(ctx context.Context) ([]User, error) {
	ctx, span := oida.Start(ctx, "SELECT users", oida.KindDatabase)
	defer span.End()

	// continue using ctx, so nested work records below this span
}
```

Pass the returned `ctx` down. That is what makes the next `Start` a child rather than a sibling.

`StartAuto` reads the name from the symbol, so it survives a rename. In the example below the span is named `storage.UserStorage.GetUsers`:

```go
ctx, span := oida.StartAuto(ctx, s.GetUsers)
defer span.End()
```

It uses reflection and the runtime symbol table, so it does not survive a stripped binary. Use `Start` with a literal name where that matters.

In a handler holding an `*http.Request`:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	r, span := oida.StartRequest(r, "user.Handler")
	defer span.End()
	// continue using r
}
```

Attributes carry the specifics, so span names stay low-cardinality:

```go
span.SetAttribute("limit", limit)
span.SetAttribute("rows", len(users))
```

## Errors

Code holding the span records on it directly:

```go
span.RecordError(err)
```

Code that does not, such as an error path several calls below, records through the context:

```go
oida.RecordError(ctx, err)
```

Both mark the span and its trace as failed, which is what the dashboard's error filter and the `Errors` column in statistics use. A nil error is ignored.

## Configuration

`platform.NewOptions` fills `Options.Telemetry` from `oida.NewOptions` and reads the environment:

| Variable                     | Default       | Meaning                                                                       |
|------------------------------|---------------|-------------------------------------------------------------------------------|
| `PLATFORM_TELEMETRY_ENABLED` | `true`        | Register the recorder and the dashboard. Disabled puts neither on the router. |
| `PLATFORM_TELEMETRY_PATH`    | `/debug/oida` | Mount path of the dashboard.                                                  |
| `PLATFORM_TELEMETRY_SERVICE` | `platform`    | Service name shown in the dashboard.                                          |

Anything else is set on the struct. See the [configuration guide](https://github.com/titpetric/oida/blob/main/docs/guide-configuration.md) for retention, sampling and access control:

```go
options := platform.NewOptions()
options.Telemetry.RingBufferSize = 500
options.Telemetry.SampleRate = 10
options.Telemetry.Authorize = func(r *http.Request) bool {
	return user.FromRequest(r).IsAdmin()
}

svc := platform.New(options)
```

The dashboard is unauthenticated unless `Authorize` says otherwise. Do not expose it publicly without one.

## Retention

Traces are kept in memory by default: a ring buffer of `RingBufferSize` traces, 200 unless you change it. It costs nothing to set up and starts empty on every boot.

```go
options := platform.NewOptions()
options.Telemetry.RingBufferSize = 500
```

To keep traces across a restart, assign disk storage. It writes one JSON document per trace into a folder, retaining at most `limit` of them:

```go
store, err := oida.NewStorageDisk(500, "/var/lib/app/traces")
if err != nil {
	return err
}

options := platform.NewOptions()
options.Telemetry.Storage = store
```

`NewStorageDisk` creates the folder and checks it is writable, so a bad path fails at startup rather than on the first recorded trace. With no path it uses a folder under the operating system temporary directory, which does not survive a reboot.

`RingBufferSize` only sizes the default memory storage. Once `Storage` is set, the `limit` argument bounds retention instead.

Storage is an interface, not a name, so it has no environment variable. A service that wants it configurable reads its own key and builds the value before calling `platform.New`.

## Bringing your own recorder

A host that registers its own telemetry module turns this one off, because two modules mounting the same path is a duplicate route and the router panics on it:

```go
options := platform.NewOptions()
options.Telemetry.Enabled = false

svc := platform.New(options)
svc.Use(recorder.Middleware)
svc.Register(recorder)
```

`Telemetry.Enabled` gates registration, not just recording: a disabled Telemetry puts no route and no middleware on the router at all. `NewTestOptions` disables it already, so tests observe no dashboard and no recorder unless they opt back in.

`PLATFORM_MODULES` does not apply. That allowlist names the application's modules, and the recorder the platform registers itself is not one of them, so it stays registered whatever the list says.

## What gets recorded

1. **Requests.** The middleware records every sampled request as a trace, named by the routed pattern, so `/users/1` and `/users/2` group into `GET /users/{id}`. The trace ID is also the `Request-Id` response header, which makes it the cheapest correlation key for logs.

2. **Startup.** Module lifecycle does not arrive over the network, so it gets a trace of its own named `platform.setup`, with `registry.Start` and one `module.start: <name>` span per module below it.

3. **Transactions.** `platform.Transaction` records a `db.Transaction` span.

Database queries are not instrumented. Spans belong at the repository boundary, where the caller knows what the query means:

```go
ctx, span := oida.Start(ctx, "SELECT users", oida.KindDatabase)
defer span.End()

rows, err := db.QueryContext(ctx, selectUsers, limit)
```

## Viewing traces

Start the service and open http://localhost:8080/debug/oida. The same paths serve JSON for tools and plain text for terminals:

```
curl -s localhost:8080/debug/oida/traces?format=json | jq '.[0].spans'
curl -s localhost:8080/debug/oida/stats?format=text
```
