package platform

import (
	"context"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/titpetric/oida"
)

// TelemetryModule records the telemetry of a platform service. It is the
// recorder and the dashboard at once: HTTP middleware recording every request,
// and a module mounting the debug front end under Options.Path.
//
// New registers one by default. A host that wires a recorder disables
// Options.Telemetry, so only one dashboard is on the path and only one
// middleware records the request.
type TelemetryModule struct {
	UnimplementedModule

	options oida.Options
	tracer  *oida.Tracer
}

var _ Module = (*TelemetryModule)(nil)

// NewTelemetryModule returns a telemetry module recording into a dedicated
// tracer, which keeps two services out of each other's traces.
func NewTelemetryModule(options oida.Options) (*TelemetryModule, error) {
	options = options.WithDefaults()
	if options.RouteFunc == nil {
		options.RouteFunc = routePattern
	}

	tracer, err := oida.New(options)
	if err != nil {
		return nil, err
	}

	return &TelemetryModule{
		UnimplementedModule: *NewUnimplementedModule("telemetry"),

		// The tracer applied the OIDA_* environment to what it was
		// handed, and the module actually runs on that copy.
		options: tracer.Options(),
		tracer:  tracer,
	}, nil
}

// Mount registers the debug front end on the platform router.
func (m *TelemetryModule) Mount(_ context.Context, r Router) error {
	return oida.Mount(r, m.tracer)
}

// Middleware records requests handled by next.
func (m *TelemetryModule) Middleware(next http.Handler) http.Handler {
	return m.tracer.Middleware(next)
}

// Options returns the options the tracer runs on, after the environment applied.
func (m *TelemetryModule) Options() oida.Options {
	return m.options
}

// Tracer returns the tracer the module records into.
func (m *TelemetryModule) Tracer() *oida.Tracer {
	return m.tracer
}

// routePattern groups statistics by the routed pattern, not the request URI.
func routePattern(r *http.Request) string {
	routeContext := chi.RouteContext(r.Context())
	if routeContext == nil {
		return ""
	}
	switch pattern := routeContext.RoutePattern(); pattern {
	case "/*", "/", "":
		return ""
	default:
		return pattern
	}
}
