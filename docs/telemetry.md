# Telemetry

The platform implements a telemetry module with OpenTelemetry.

It is included by importing `github.com/titpetric/platform/module/telemetry`.
The package provides a middleware and instrumentation for the database.

```
var TracerProvider *sdktrace.TracerProvider
    TracerProvider is shared globally


FUNCTIONS

func CaptureError(ctx context.Context, err error)
    CaptureError logs an error that occured in a request.

func Connect(driver, dsn string) (*sqlx.DB, error)
    Connect is like Open but verifies the connection (calls Ping).

func Middleware(next http.Handler) http.Handler
    Middleware returns an otelhttp middleware

func Open(driver, dsn string) (*sqlx.DB, error)
    Open creates a new instrumented *sqlx.DB without pinging. It wraps
    otelsql.Open and sqlx.NewDb for full compatibility.
```

The package also implements an `init` function to set up the
opentelemetry collector. To test the platform with opentelemetry, `task
docker up` in the root of the project will build the app and bring up
the required services.

To access the dashboard, open: [http://localhost:16686/](http://localhost:16686/).

Requests made to the service will be logged in opentelemetry.

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
      - OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318
```

Check the rest of the `docker-compose.yml` in the repository to see the
setup for other observability components.

### Tasks

- `task docker` - builds the platform Docker image.  
- `task up` - starts the Docker test environment. The database is external to the platform container.
- `task down` - stops the Docker test environment.

### Observability Features

1. **APM Tracing**:
   Requests to the platform app are automatically traced in Elastic APM. Chi route middleware is instrumented by default, so HTTP spans will be visible in the APM dashboard.

2. **Database Instrumentation**:
   Supported drivers (SQLite, MySQL) are instrumented. Queries executed via the platform modules are automatically captured as spans.

3. **Explicit Error Capture**:
   Any error you want to log to the observability platform can be captured explicitly:

```go
telemetry.CaptureError(ctx, err)
```

Use this in module handlers or background tasks to report errors that
occur during request processing. Background tasks need to create their
own transaction context.
