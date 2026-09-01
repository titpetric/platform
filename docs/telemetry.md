# Telemetry

The platform records telemetry with [oida](https://github.com/titpetric/oida): traces and spans held in a ring buffer inside the process, with a dashboard served at `/debug/oida`. There is no collector, no exporter and no second service to run.

Recording is opt-in. `platform.New` registers the recorder and the dashboard only when `Options.Telemetry.Enabled` is set, because the dashboard reports the internals of the process and is unauthenticated until it is configured otherwise. Ask for it with the environment:

```bash
PLATFORM_TELEMETRY_ENABLED=true
```

or in code:

```go
options := platform.NewOptions()
options.Telemetry.Enabled = true

svc := platform.New(options)
```

With telemetry off, the instrumentation calls below still compile and run: they record nothing and cost a nil check.

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

## Logs

A trace carries log entries, attributed to the innermost open span, so the lines a request wrote are read next to the spans that wrote them:

```go
span.Info("cache miss", "key", key)
span.Error("upstream refused", "status", resp.StatusCode)
```

`CaptureLogs` is on by default, and `OIDA_CAPTURE_LOGS=false` turns it off. Disabled, `Info` does nothing and `Error` records its formatted text on the active span the way `RecordError` does, so the message is not lost.

## Configuration

`platform.NewOptions` fills `Options.Telemetry` from `oida.NewOptions` and reads the environment:

| Variable                     | Default       | Meaning                                                                  |
|------------------------------|---------------|--------------------------------------------------------------------------|
| `PLATFORM_TELEMETRY_ENABLED` | `false`       | Register the recorder and the dashboard. Off puts neither on the router. |
| `PLATFORM_TELEMETRY_PATH`    | `/debug/oida` | Mount path of the dashboard.                                             |
| `PLATFORM_TELEMETRY_SERVICE` | `platform`    | Service name shown in the dashboard.                                     |

`oida.NewOptions` also sets `ReadEnv`, so the recorder applies its own `OIDA_*` variables when the tracer is built: retention, sampling, allowed networks, users and the signing secret are configurable from the deployment without the platform proxying a variable for each of them. They are the table in the [configuration guide](https://github.com/titpetric/oida/blob/main/docs/guide-configuration.md), and an `OIDA_*` variable applies only where the code left the field at its default, so a value set on `Options.Telemetry` wins.

The three variables above are the exception, because they are read before a tracer exists. `PLATFORM_TELEMETRY_SERVICE` and `PLATFORM_TELEMETRY_PATH` are set on the struct, so `OIDA_SERVICE_NAME` never applies and `OIDA_PATH` applies only when the platform variable is unset. `OIDA_ENABLED` is read by the platform alongside its own variable, because registration is decided before the tracer that would have applied it; `PLATFORM_TELEMETRY_ENABLED` decides when both are set.

Anything else is set on the struct:

```go
options := platform.NewOptions()
options.Telemetry.RingBufferSize = 500
options.Telemetry.SampleRate = 10
options.Telemetry.Authorize = func(r *http.Request) bool {
	return user.FromRequest(r).IsAdmin()
}

svc := platform.New(options)
```

The dashboard is unauthenticated unless it is told otherwise. `Authorize` is the platform's own hook; the recorder adds a CIDR allow list, a login screen and bearer tokens:

```go
options.Telemetry.AllowedNetworks = []string{"127.0.0.0/8", "10.0.0.0/8"}
options.Telemetry.Users = map[string]string{"admin": bcryptHash} // htpasswd -nbB admin secret
options.Telemetry.SigningSecret = os.Getenv("OIDA_SIGNING_SECRET")
```

`Authorize` runs first, then the allow list, then credentials. A request rejected by either of the first two gets a 404; one missing credentials is redirected to `{path}/login` or answered with a 401. Do not expose the dashboard publicly without at least one of them.

## Retention

Traces are kept in memory by default: a ring buffer of `RingBufferSize` traces, 200 unless you change it. It costs nothing to set up and starts empty on every boot.

```go
options := platform.NewOptions()
options.Telemetry.RingBufferSize = 500
```

To keep traces across a restart, ask for disk storage. It writes one JSON document per trace into a folder, retaining at most the limit it is given. The recorder builds it from the environment:

```bash
OIDA_STORAGE_DRIVER=disk
OIDA_STORAGE_DISK_PATH=/var/lib/app/traces
OIDA_STORAGE_DISK_LIMIT=5000
```

or from code, with the driver from `github.com/titpetric/oida/storage`:

```go
store, err := storage.NewDiskStorage(500, "/var/lib/app/traces")
if err != nil {
	return err
}

options := platform.NewOptions()
options.Telemetry.Storage = store
```

Either way the folder is created and checked for writability, so a bad path fails at startup rather than on the first recorded trace. With no path it uses a folder under the operating system temporary directory, which does not survive a reboot.

`RingBufferSize` only sizes the default memory storage. Once `Storage` is set, the driver's own limit bounds retention instead, and it wins over every `OIDA_STORAGE_*` variable.

## Bringing your own recorder

A host that registers its own telemetry module turns this one off, so only one dashboard is on the path and only one middleware records the request:

```go
options := platform.NewOptions()
options.Telemetry.Enabled = false

svc := platform.New(options)
svc.Use(recorder.Middleware)
svc.Register(recorder)
```

`Telemetry.Enabled` gates registration, not just recording: a disabled Telemetry puts no route and no middleware on the router at all. `NewTestOptions` disables it already, so tests observe no dashboard and no recorder unless they opt back in.

`Options.Telemetry` is `oida.Options`, so the zero value is disabled. An `Options` built by hand asks for the recorder the same way `NewOptions` does, by filling it from `oida.NewOptions` and enabling it:

```go
options := &platform.Options{
	ServerAddr: ":8080",
	Telemetry:  oida.NewOptions("billing-api"),
}
options.Telemetry.Enabled = true
```

`oida.NewOptions` also sets `ReadEnv`. A recorder that must ignore the deployment, which is what a test wants, clears it:

```go
options.Telemetry.ReadEnv = false
```

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
