package platform_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/titpetric/oida"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

// telemetryPath is the path these tests mount the dashboard on. It is spelled
// out rather than read from the recorder, so a test asks for the address it
// configured instead of whatever the library defaults to.
const telemetryPath = "/debug/oida"

// tracedOptions returns recorder options a test can rely on: recording is on,
// memory accounting is off because it costs a GC read per trace, and ReadEnv is
// off so an OIDA_* variable in the environment cannot reconfigure the run.
func tracedOptions() oida.Options {
	options := oida.NewOptions("platform-test")
	options.Path = telemetryPath
	options.Enabled = true
	options.TrackMemoryUse = false
	options.ReadEnv = false
	return options
}

// newTracedPlatform returns a platform recording into oida. NewTestOptions
// disables telemetry, so a test that wants the dashboard opts back in.
func newTracedPlatform(t *testing.T) *platform.Platform {
	t.Helper()

	options := platform.NewTestOptions()
	options.Telemetry = tracedOptions()

	svc := platform.New(options)
	t.Cleanup(svc.Stop)

	svc.Register(&platform.UnimplementedModule{
		NameFn: func() string { return "TestTelemetry" },
		MountFn: func(_ context.Context, r platform.Router) error {
			r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
				_, span := oida.Start(r.Context(), "load user", oida.KindDatabase)
				defer span.End()

				_, err := w.Write([]byte("user"))
				require.NoError(t, err)
			})
			return nil
		},
	})

	require.NoError(t, svc.Start(t.Context()))
	return svc
}

func get(t *testing.T, url string) (int, []byte) {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, resp.Body.Close()) })

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, body
}

func TestTelemetryModuleServesDashboard(t *testing.T) {
	svc := newTracedPlatform(t)

	status, _ := get(t, svc.URL()+"/users/42")
	require.Equal(t, http.StatusOK, status)

	status, body := get(t, svc.URL()+telemetryPath)
	require.Equal(t, http.StatusOK, status)
	require.True(t, len(body) > 0)
}

func TestTelemetryModuleRecordsRequests(t *testing.T) {
	svc := newTracedPlatform(t)

	status, _ := get(t, svc.URL()+"/users/42")
	require.Equal(t, http.StatusOK, status)

	status, body := get(t, svc.URL()+telemetryPath+"/traces?format=json")
	require.Equal(t, http.StatusOK, status)

	var traces []oida.Trace
	require.NoError(t, json.Unmarshal(body, &traces))

	// The startup trace and the request. Requests to the dashboard itself
	// are not recorded, so the count is stable.
	byName := map[string]oida.Trace{}
	for _, trace := range traces {
		byName[trace.Name] = trace
	}

	request, ok := byName["GET /users/{id}"]
	require.True(t, ok)
	require.Equal(t, http.StatusOK, request.HTTP.Status)

	// RouteFunc groups by the routed pattern, not the request URI, which is
	// what keeps /users/1 and /users/2 in one statistics row.
	var names []string
	for _, span := range request.Spans {
		names = append(names, span.Name)
	}
	require.Contains(t, names, "load user")
}

func TestTelemetryModuleRecordsStartup(t *testing.T) {
	svc := newTracedPlatform(t)

	status, body := get(t, svc.URL()+telemetryPath+"/traces?format=json")
	require.Equal(t, http.StatusOK, status)

	var traces []oida.Trace
	require.NoError(t, json.Unmarshal(body, &traces))

	var setup *oida.Trace
	for i, trace := range traces {
		if trace.Name == "platform.setup" {
			setup = &traces[i]
			break
		}
	}
	require.NotNil(t, setup)

	// Module lifecycle spans nest below the startup trace. Without the
	// background trace opened in setup they would have nowhere to land.
	var names []string
	for _, span := range setup.Spans {
		names = append(names, span.Name)
	}
	require.Contains(t, names, "registry.Start")
	require.Contains(t, names, "module.start: TestTelemetry")
}

// A disabled Telemetry registers no module, so the dashboard is not on the
// router at all rather than mounted and empty.
func TestTelemetryDisabledRegistersNoModule(t *testing.T) {
	svc := NewTestPlatform(t)

	status, _ := get(t, svc.URL()+telemetryPath)
	require.Equal(t, http.StatusNotFound, status)
}

func TestTelemetryModuleSurvivesModuleFilter(t *testing.T) {
	options := platform.NewTestOptions()
	options.Telemetry = tracedOptions()

	// The allowlist names the application's modules. The recorder the
	// platform registers itself is not one of them.
	options.Modules = []string{"TestTelemetry"}

	svc := platform.New(options)
	t.Cleanup(svc.Stop)
	svc.Register(&platform.UnimplementedModule{
		NameFn:  func() string { return "TestTelemetry" },
		MountFn: func(context.Context, platform.Router) error { return nil },
	})
	require.NoError(t, svc.Start(t.Context()))

	status, _ := get(t, svc.URL()+telemetryPath)
	require.Equal(t, http.StatusOK, status)
}

// Recording is opt-in: off in the zero value, and off in the options both
// constructors hand out.
func TestTelemetryDefaults(t *testing.T) {
	require.False(t, platform.NewOptions().Telemetry.Enabled)
	require.False(t, platform.NewTestOptions().Telemetry.Enabled)

	// A service built from the defaults serves no dashboard, so nothing
	// reports the internals of a process that never asked to.
	options := platform.NewOptions()
	options.ServerAddr, options.Quiet = "127.0.0.1:0", true
	svc := platform.New(options)
	t.Cleanup(svc.Stop)
	require.NoError(t, svc.Start(t.Context()))

	status, _ := get(t, svc.URL()+options.Telemetry.Path)
	require.Equal(t, http.StatusNotFound, status)
}

// The environment turns recording on, which is how a deployment asks for the
// dashboard without a code change.
func TestTelemetryEnabledByEnv(t *testing.T) {
	t.Setenv("PLATFORM_TELEMETRY_ENABLED", "true")

	options := platform.NewOptions()
	require.True(t, options.Telemetry.Enabled)

	options.ServerAddr, options.Quiet = "127.0.0.1:0", true
	options.Telemetry.TrackMemoryUse = false
	svc := platform.New(options)
	t.Cleanup(svc.Stop)
	require.NoError(t, svc.Start(t.Context()))

	status, _ := get(t, svc.URL()+options.Telemetry.Path)
	require.Equal(t, http.StatusOK, status)
}

// OIDA_ENABLED is read by the platform rather than left to oida, because
// registration is decided before a tracer exists: a module that never
// registered never reaches the New that would have applied the variable.
func TestTelemetryEnabledByRecorderEnv(t *testing.T) {
	// The test environment pins the platform's own variable, so it is the
	// fallback that is under test here.
	t.Setenv("PLATFORM_TELEMETRY_ENABLED", "")
	t.Setenv("OIDA_ENABLED", "true")
	require.True(t, platform.NewOptions().Telemetry.Enabled)

	// The platform's own variable decides, so a deployment can keep the
	// dashboard off a service the recorder is configured for elsewhere.
	t.Setenv("PLATFORM_TELEMETRY_ENABLED", "false")
	require.False(t, platform.NewOptions().Telemetry.Enabled)
}
