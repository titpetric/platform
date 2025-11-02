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
	if os.Getenv("PLATFORM_ENABLE_EXPVAR") == "true" {
		monitor.enabled = true
	}

	if os.Getenv("PLATFORM_ENABLE_OTEL") != "true" {
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
}

// Start is a wrapper to tracer.Start and expvar.NewInt.
//
// Calling the function adds a span to the current context. The name should
// be a symbol reference for the caller. For best results, combine the
// package, symbol and function into a name like `user.service.Login`.
func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	monitorTouch(name)
	return tracer.Start(ctx, name)
}

// StartAuto tries to fill the span name from the symbol.
//
// It's intended to pass a function, or a type. The package name, type
// name, and function name are combined with `.` to delimit them.
// See tests under internal/reflect for more information.
//
// StartAuto uses reflection and may not work correctly under various conditions.
// If performance or build restrictions are impacting use, use `Start`.
func StartAuto(ctx context.Context, symbol any) (context.Context, trace.Span) {
	name := reflect.SymbolName(symbol)
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

// Fatal wraps CaptureError with log.Fatal.
func Fatal(ctx context.Context, err error) {
	CaptureError(ctx, err)
	// Do we need span.Send?
	log.Fatal(err)
}
