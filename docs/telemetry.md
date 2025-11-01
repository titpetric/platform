# Telemetry

The platform implements a telemetry module with OpenTelemetry.

To test the platform with OTel, `task docker up` in the root
of the project will build the app and bring up the required services.
Check the `docker-compose.yml` file for configuration.

To access the tracing dashboard, open: [http://localhost:16686/](http://localhost:16686/).

## Usage

Instrumenting your code with telemetry is simple. To use telemetry,
start by importing `github.com/titpetric/platform/telemetry`. This
provides functions to use `context.Context` and `*http.Request` values.

```
func Start(ctx context.Context, name string) (context.Context, trace.Span)
    Start is a wrapper to tracer.Start. It's meant to add instrumentation in the
    storage layer, or around important bits of code. It adds nothing to the span
    but the name. Ideally use a FQDN ("package.Type.Function").

func StartRequest(r *http.Request, name string) (*http.Request, trace.Span)
    StartRequest is an utility to take the http.Request and update it's context.

func StartAuto(ctx context.Context, symbol any) (context.Context, trace.Span)
    StartAuto tries to fill the span name from the symbol.

    It's intended to pass a function, or a type. The package name, type name,
    and function name are combined with `.` to delimit them. See tests under
    internal/reflect for more information.
```

The package exposes other symbols, but don't rely on them.

So, to trace from a `*http.Request`:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	r, span = telemetry.StartRequest(r, "user.Handler")
	defer span.End()
	// continue using `r`
}
```

And to trace any context aware function:

```go
func (*UserStorage) GetUsers(ctx context.Context) []string {
	ctx, span = telemetry.Start(ctx, "storage.UserStorage.GetUsers")
	defer span.End()

	// continue using ctx for database queries
}
```

Or with `StartAuto`, passing the symbol to trace by name. In the given
example, the name of the symbol is `storage.UserStorage.GetUsers`.

```go
func (*UserStorage) GetUsers(ctx context.Context) []string {
	ctx, span = telemetry.StartAuto(ctx, GetUsers)
	defer span.End()

	// continue using ctx for database queries
}
```

> StartAuto uses reflection and may not work correctly under various conditions.
> If performance or build restrictions are impacting use, use `Start`.

The `span` value notably contains an implementation of the `trace.Span`
interface. In that interface is a `SetName(string)` function that lets
you customize the name of the span.

For example, platform requests carry the request method and the route
pattern, e.g. `POST /login`, to ease the grouping process. It's nice to
see repeated traces, but the behaviour depends a lot on the tool used to
display opentelemetry data. This seems suitable for Jaeger and Elastic
APM.

## Docker Test Environment

You can quickly bring up a test environment with Docker Compose. Example configuration:

```yaml
  platform:
    image: titpetric/platform
    networks:
      - elastic
      - db_storage
    ports:
      - 8080:8080
    environment:
      - PLATFORM_DB_DEFAULT=mysql://platform:platform@tcp(db1:3306)/platform
      - OTEL_SERVICE_NAME=platform
      - OTEL_SERVICE_ENABLED=true
      - OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318
```

Check the rest of the `docker-compose.yml` in the repository to see the
setup for other observability components.

### Tasks

- `task docker` - builds the platform Docker image.  
- `task up` - starts the Docker test environment. The database is external to the platform container.
- `task down` - stops the Docker test environment.

### Observability Features

The OpenTelemetry stack includes their collector service, prometheus and jaeger.

1. **APM Tracing**: Requests to the platform app are automatically traced in Jaeger. Other OpenTelemetry monitoring tooling exists.

2. **Database Instrumentation**: Any database driver in use is instrumented with telemetry. 

   Supported drivers (SQLite, MySQL) are instrumented. Queries executed via the platform modules are automatically captured as spans.

3. **Explicit Error Capture**:
   Any error you want to log to the observability platform can be captured explicitly:

```go
telemetry.CaptureError(ctx, err)
```

Use this in module handlers or background tasks to report errors that
occur during request processing. Background tasks need to create their
own transaction context.
