package platform

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/titpetric/platform/pkg/httpcontext"
)

// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr string

	// Quiet turns down the verbosity in the Platform logging code, set to true in tests.
	Quiet bool

	// Modules controls which modules get loaded. If the list
	// is empty (unconfigured, zero value), all modules load.
	Modules []string

	// ConfigFS can be used for configuration purposes by modules. It's optional and may be nil.
	// The application running with the platform may use `go:embed` to carry config for the
	// composed service.
	ConfigFS fs.FS

	// Telemetry configures the recorder and the debug dashboard. A
	// disabled Telemetry registers no module at all.
	Telemetry Telemetry
}

// NewOptions provides default options for the platform.
func NewOptions() *Options {
	opt := &Options{}
	opt.ServerAddr = opt.env("PLATFORM_SERVER_ADDR", ":8080")
	opt.Modules = opt.envCSV("PLATFORM_MODULES")

	opt.Telemetry = NewTelemetry()
	opt.Telemetry.ServiceName = opt.env("PLATFORM_TELEMETRY_SERVICE", "platform")
	opt.Telemetry.Path = opt.env("PLATFORM_TELEMETRY_PATH", opt.Telemetry.Path)
	opt.Telemetry.Enabled = opt.envBool("PLATFORM_TELEMETRY_ENABLED", opt.Telemetry.Enabled)
	return opt
}

func (*Options) envCSV(name string) []string {
	if v := os.Getenv(name); v != "" {
		return strings.Split(v, ",")
	}
	return nil
}

func (*Options) env(name string, def string) string {
	result := def
	if v := os.Getenv(name); v != "" {
		result = v
	}
	return result
}

func (*Options) envBool(name string, def bool) bool {
	v, err := strconv.ParseBool(os.Getenv(name))
	if err != nil {
		return def
	}
	return v
}

// NewTestOptions produces default options for tests.
func NewTestOptions() *Options {
	return &Options{
		ServerAddr: "127.0.0.1:0",
		Quiet:      true,

		// Tests get no dashboard and no recorder unless they ask for
		// one, which keeps the global state they observe empty.
		Telemetry: Telemetry{},
	}
}

type optionsKey struct{}

var optionsContext = httpcontext.NewValue[*Options](optionsKey{})

// OptionsFromRequest returns the *Options instance attached to the request.
func OptionsFromRequest(r *http.Request) *Options {
	return optionsContext.Get(r)
}

// OptionsFromContext returns the *Options instance attached to the context.
func OptionsFromContext(ctx context.Context) *Options {
	return optionsContext.GetContext(ctx)
}
