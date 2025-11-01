package telemetry

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/XSAM/otelsql"
	"github.com/riandyrn/otelchi"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/titpetric/platform/internal/reflect"
)

var tracer trace.Tracer

const Name = "internal/telemetry"

func init() {
	if os.Getenv("OTEL_SERVICE_ENABLED") != "true" {
		log.Println("[telemetry] OpenTelemetry is disabled, setting tracer to noop")
		tracer = trace.NewNoopTracerProvider().Tracer(Name)
		return
	}

	initOpenTelemetry()
}

func initOpenTelemetry() {
	// instrument sql
	sql_open = func(driver, dsn string) (*sql.DB, error) {
		return otelsql.Open(driver, dsn)
	}

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4318" // default local collector
	}

	ctx := context.Background()

	// Create OTLP trace exporter via HTTP
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("telemetry: failed to create OTLP exporter: %v", err)
	}

	// Define the Resource (service name, version, etc.)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("platform"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("telemetry: failed to create resource: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	tracer = tracerProvider.Tracer(Name)

	otel.SetTracerProvider(tracerProvider)

	log.Println("[telemetry] OpenTelemetry initialized")
}

// Start is a wrapper to tracer.Start. It's meant to add instrumentation
// in the storage layer, or around important bits of code. It adds nothing
// to the span but the name. Ideally use a FQDN ("package.Type.Function").
func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

// StartAuto tries to fill the span name from the symbol. If a pointer of
// a type is passed, the expected output contains the package and the type
// name concatenated with a dot. If a function is passed, the function
// name is appended with another dot. For example:
//
// - internal.StartAuto
// - storage.UserStorage.Get
//
// It's good enough considering handlers should carry parent span context.
// If a string is passed, the string is returned as is.
//
// In addition, the name will be used to publish an expvar counter.
func StartAuto(ctx context.Context, symbol any) (context.Context, trace.Span) {
	name := reflect.SymbolName(symbol)
	monitorTouch(name)
	return Start(ctx, name)
}

// StartRequest is an utility to take the http.Request and update it's context.
func StartRequest(r *http.Request, name string) (*http.Request, trace.Span) {
	ctx := r.Context()
	ctx, span := Start(ctx, name)
	return r.WithContext(ctx), span
}

// Middleware returns a middleware that instruments requests with telemetry.
func Middleware(name string) func(http.Handler) http.Handler {
	return otelchi.Middleware(name, otelchi.WithRequestMethodInSpanName(true))
}

// CaptureError logs an error that occured in a request.
func CaptureError(ctx context.Context, err error) {
	trace.SpanFromContext(ctx).RecordError(err)
}
