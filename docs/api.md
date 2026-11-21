# Package platform

```go
import (
	"github.com/titpetric/platform"
)
```

The platform is an extensible modular system for writing HTTP servers.

1. Provides a global registry for middleware and module registration
2. Provides a lifecycle to the modules for graceful shutdown
3. Provides a router the modules can attach to

It's advised to use `platform.RegisterFunc` from `init` functions. Similarly, `platform.Use` belongs to `main` at the top level, or to any descendant setup function. Don't use these functions from tests as they create a shared state.

It's possible to use the platform in an imperative way.

```go
svc := platform.New(platform.NewOptions())
svc.Use(middleware.Logger)
svc.Register(user.NewModule())
```

The platform lifecycle is extensively tested to ensure no races, no goroutine leaks. Each platform object creates a copy of the global state and holds scoped allocations only, enabling test parallelism. Modules are part of that copy when registered with `RegisterFunc`, which is called once per platform. A value registered with the deprecated `Register` is shared by the platforms in the process.

## Types

<details>
<summary><code>type DatabaseProvider</code></summary>

```go
// DatabaseProvider opens named connections, "default" when no name is given.
type DatabaseProvider interface {
	// Open takes a context, a different implementation may use it for context awareness.
	// For example, Open may retrieve connection details from the database. The context
	// is used for tracing that operation, and the error returned may be the
	// result of a query issued against the database.
	Open(ctx context.Context, names ...string) (*sql.DB, error)

	// Connect has the semantics of Open + PingContext against the database,
	// verifying that the connection is live. An error is returned if
	// the storage is unreachable.
	//
	// A real database provider options may dictate additional behaviour.
	// For example, if a connection fails, it may retry a number of times,
	// it may back-off, it may continue to retry the connection in a
	// blocking way. It may also grab additional fail-over information from
	// keys based on the input parameters. An example of that would be to
	// define `name` and `name/failover` connections. If the connection to
	// the first fails, the second upstream is used.
	Connect(ctx context.Context, names ...string) (*sql.DB, error)
}
```

</details>

<details>
<summary><code>type ErrorResponse</code></summary>

```go
// ErrorResponse is the JSON structure for error responses.
type ErrorResponse struct {
	Error ErrorResponseBody `json:"error"`
}
```

</details>

<details>
<summary><code>type ErrorResponseBody</code></summary>

```go
// ErrorResponseBody is an inner type for ErrorResponse.Error.
type ErrorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
```

</details>

<details>
<summary><code>type Logger</code></summary>

```go
// Logger is the interface the platform writes output through: the subset of
// *slog.Logger it needs. A *slog.Logger assigns to Platform.Logger as it is.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
```

</details>

<details>
<summary><code>type Manager</code></summary>

```go
// Manager owns the socket and the platform generation serving there. A SIGHUP
// replaces the generation and keeps the address.
type Manager struct {
	// Logger receives the manager's own output, and is handed to every
	// platform generation it creates. Defaults as Platform.Logger does.
	Logger Logger

	// Setup runs against every platform generation before it starts.
	// Registration against the platform value belongs here, as a reload
	// discards the value it was made against. Assign it ahead of Start: a
	// reload runs it on the signal handler’s goroutine.
	Setup func(*Platform) error

	// Check runs before retiring the serving generation, and a
	// non-nil error abandons that reload with the generation left alone.
	// Reading the configuration a reload would apply belongs here: Setup
	// runs against the new generation, which exists only once the old one
	// has stopped. Nil reloads unconditionally. Assign before Start, as
	// with Setup, because the signal handler reads it on a goroutine.
	Check func() error

	options *Options

	// mu serializes the generation swap.
	mu sync.Mutex

	// shared is the socket every generation serves, set once by Start. It is
	// read outside mu by URL, which must not wait out a reload.
	shared  atomic.Pointer[sharedListener]
	current atomic.Pointer[generation]

	// pid records this process's id. The manager holds it rather than the
	// platform, because a reload replaces the platform and not the process.
	pid pidfile.Pidfile

	// final shutdown context, cancelled when the manager stops
	context context.Context
	cancel  context.CancelFunc
	stop    func()
	once    sync.Once
}
```

</details>

<details>
<summary><code>type Middleware</code></summary>

```go
// Middleware is a type alias for middleware functions.
type Middleware func(http.Handler) http.Handler
```

</details>

<details>
<summary><code>type Module</code></summary>

```go
// Module is the implementation contract for modules.
// The interface should only be used to enforce the API contract as
// shown below. It's also used to provide `platform.Register()`.
type Module interface {
	// Name should return a meaningful name for your module.
	Name() string

	// Start is used to create any goroutines or otherwise
	// set up the module by starting a server. It allows
	// to implement a lifecycle of the service.
	Start(context.Context) error

	// Stop should clean up any goroutines, clean up leaks.
	Stop(context.Context) error

	// Mount runs before the server starts, and allows you to
	// register new routes to your module.
	Mount(context.Context, Router) error
}
```

</details>

<details>
<summary><code>type Options</code></summary>

```go
// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr string

	// PidFile is the file the process records its id in, for a service
	// manager or a command line that signals it. Empty writes none. The
	// directory must exist, an existing file is overwritten, and the file
	// is removed on a clean stop. A Manager writes it for the process, and the
	// reload leaves it alone.
	PidFile string

	// Quiet silences the platform's own output: New installs a discarding
	// logger as Platform.Logger instead of the default one. Set to true in
	// tests. Assigning Platform.Logger afterwards overrules it.
	Quiet bool

	// Modules controls which modules get loaded. If the list
	// is empty (unconfigured, zero value), all modules load.
	Modules []string

	// ConfigFS can be used for configuration purposes by modules. It's optional and may be nil.
	// The application running with the platform may use `go:embed` to carry config for the
	// composed service.
	ConfigFS fs.FS

	// Telemetry configures the recorder and the debug dashboard. It is
	// off unless asked for: the zero value disables it, and NewOptions
	// disables it too, because the dashboard reports the internals of the
	// process and is unauthenticated unless Telemetry.Authorize, its
	// network allow list or its users say otherwise. Set Enabled, or
	// PLATFORM_TELEMETRY_ENABLED, to record.
	Telemetry oida.Options
}
```

</details>

<details>
<summary><code>type Platform</code></summary>

```go
// Platform is our world struct.
type Platform struct {
	// Logger receives the platform's own output. New sets it to
	// slog.Default(), or to a discarding logger when Options.Quiet is set.
	// Assign to it ahead of Start to reuse the platform logger from a test,
	// or to route output into a consumer application's logger.
	Logger Logger

	options *Options

	// server setup
	router   *chi.Mux
	server   *http.Server
	listener net.Listener

	// served is closed when the accept loop has returned. Stop waits for the socket,
	// and the manager gets the socket back when Stop is done.
	served chan struct{}

	// final shutdown context
	context  context.Context
	cancel   context.CancelFunc
	stop     func()
	once     sync.Once
	stopping atomic.Bool

	// managed marks a generation a Manager runs. It keeps the generation
	// from arming the process signals, which the manager owns.
	managed bool

	// pid records this process's id, disabled when Options.PidFile is
	// empty. A Manager clears it on the generations it runs and holds its
	// own, because the file records a process and a reload makes no new one.
	pid pidfile.Pidfile

	// registry holds settings for plugins and middleware.
	// It's auto-filled from global scope.
	registry *Registry

	// telemetry records requests and serves the debug dashboard. It's nil
	// when Options.Telemetry is disabled, or when the recorder failed to
	// build, in which case instrumentation degrades to no-ops.
	telemetry *TelemetryModule
}
```

</details>

<details>
<summary><code>type Registry</code></summary>

```go
// Registry provides a programmatic API to manage middleware and modules.
// A module registers middleware and has a contract to enforce lifecycle.
type Registry struct {
	mu sync.RWMutex

	// Registrations hold every module registered, in registration order.
	// The list is filtered to start/stop only the modules that are enabled.
	// Interacting with a module is subject to concurrency concerns.
	registrations []registration
	middleware    []Middleware

	// On registry start when modules start, a cleanup service per module
	// will be registered via this value. On registry close, the slice
	// is cleared. The functions receive the shutdown context.
	cleanups []func(context.Context)
}
```

</details>

<details>
<summary><code>type Router</code></summary>

```go
// Router is a local shim that aliases the chi router interface.
type Router = chi.Router
```

</details>

<details>
<summary><code>type TelemetryModule</code></summary>

```go
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
```

</details>

<details>
<summary><code>type UnimplementedModule</code></summary>

```go
// UnimplementedModule implements the module contract.
// The module can embed the type to skip implementing
// any of the bound functions.
type UnimplementedModule struct {
	NameFn  func() string
	StartFn func(context.Context) error
	StopFn  func(context.Context) error
	MountFn func(context.Context, Router) error
}
```

</details>

## Vars

<details>
<summary><code>var Database</code></summary>

```go
// Database is a holder of the database provider api in package namespace.
var Database DatabaseProvider = global.db
```

</details>

## Function symbols

- `func Error (w http.ResponseWriter, r *http.Request, status int, data error)`
- `func FromContext (ctx context.Context) *Platform`
- `func FromRequest (r *http.Request) *Platform`
- `func JSON (w http.ResponseWriter, r *http.Request, status int, data any)`
- `func New (options *Options) *Platform`
- `func NewManager (options *Options) *Manager`
- `func NewOptions () *Options`
- `func NewTelemetryModule (options oida.Options) (*TelemetryModule, error)`
- `func NewTestOptions () *Options`
- `func NewUnimplementedModule (name string) *UnimplementedModule`
- `func OptionsFromContext (ctx context.Context) *Options`
- `func OptionsFromRequest (r *http.Request) *Options`
- `func Param (r *http.Request, name string) string`
- `func QueryParam (r *http.Request, name string) string`
- `func Register (m Module)`
- `func RegisterFunc (f func() Module)`
- `func SetupConnections (environment []string)`
- `func Start (ctx context.Context, options *Options) (*Platform, error)`
- `func TestMiddleware () Middleware`
- `func Transaction (ctx context.Context, db *sql.DB, fn func(context.Context, *sql.Tx) error) error`
- `func URLParam (r *http.Request, name string) string`
- `func Use (mw Middleware)`
- `func (*Manager) Context () context.Context`
- `func (*Manager) Platform () *Platform`
- `func (*Manager) Reload (ctx context.Context) error`
- `func (*Manager) Start (ctx context.Context) error`
- `func (*Manager) Stop ()`
- `func (*Manager) URL () string`
- `func (*Manager) Wait ()`
- `func (*Platform) Context () context.Context`
- `func (*Platform) Find (target any) bool`
- `func (*Platform) Register (m Module)`
- `func (*Platform) Start (ctx context.Context) error`
- `func (*Platform) Stats () (int, int)`
- `func (*Platform) Stop ()`
- `func (*Platform) URL () string`
- `func (*Platform) Use (m Middleware)`
- `func (*Platform) Wait ()`
- `func (*Registry) Cleanup (fn func(context.Context))`
- `func (*Registry) Clone () *Registry`
- `func (*Registry) Close (ctx context.Context)`
- `func (*Registry) Find (target any) bool`
- `func (*Registry) Register (m Module)`
- `func (*Registry) RegisterFunc (f func() Module)`
- `func (*Registry) Start (ctx context.Context, mux Router, opts *Options) error`
- `func (*Registry) Stats () int`
- `func (*Registry) Use (f Middleware)`
- `func (*TelemetryModule) Middleware (next http.Handler) http.Handler`
- `func (*TelemetryModule) Mount (_ context.Context, r Router) error`
- `func (*TelemetryModule) Options () oida.Options`
- `func (*TelemetryModule) Tracer () *oida.Tracer`
- `func (UnimplementedModule) Mount (ctx context.Context, r Router) error`
- `func (UnimplementedModule) Name () string`
- `func (UnimplementedModule) Start (ctx context.Context) error`
- `func (UnimplementedModule) Stop (ctx context.Context) error`

### `Error`

`Error` writes an error payload as JSON.

```go
func Error(w http.ResponseWriter, r *http.Request, status int, data error)
```

### `FromContext`

`FromContext` returns the *Platform instance attached to the context.

```go
func FromContext(ctx context.Context) *Platform
```

### `FromRequest`

`FromRequest` returns the *Platform instance attached to the request.

```go
func FromRequest(r *http.Request) *Platform
```

### `JSON`

`JSON` writes any payload as JSON. If the payload is nil, the write is omitted. If an error occurs in encoding, a telemetry error is logged.

```go
func JSON(w http.ResponseWriter, r *http.Request, status int, data any)
```

### `New`

`New` will create a new *Platform object. It is the allocation point for each platform instance. If no options are passed, the defaults are in use. The defaults options are provided by NewOptions().

```go
func New(options *Options) *Platform
```

### `NewManager`

`NewManager` creates a manager for the passed options. If no options are passed, the defaults from NewOptions() are in use.

```go
func NewManager(options *Options) *Manager
```

### `NewOptions`

`NewOptions` provides default options for the platform.

```go
func NewOptions() *Options
```

### `NewTelemetryModule`

`NewTelemetryModule` returns a telemetry module recording into a dedicated tracer, which keeps two services out of each other's traces.

```go
func NewTelemetryModule(options oida.Options) (*TelemetryModule, error)
```

### `NewTestOptions`

`NewTestOptions` produces default options for tests.

```go
func NewTestOptions() *Options
```

### `NewUnimplementedModule`

`NewUnimplementedModule` will fill the module name.

```go
func NewUnimplementedModule(name string) *UnimplementedModule
```

### `OptionsFromContext`

`OptionsFromContext` returns the *Options instance attached to the context.

```go
func OptionsFromContext(ctx context.Context) *Options
```

### `OptionsFromRequest`

`OptionsFromRequest` returns the *Options instance attached to the request.

```go
func OptionsFromRequest(r *http.Request) *Options
```

### `Param`

`Param` will return a named URL parameter, or query string.

```go
func Param(r *http.Request, name string) string
```

### `QueryParam`

`QueryParam` will return a named query parameter from the request.

```go
func QueryParam(r *http.Request, name string) string
```

### `Register`

`Register` adds a module to the global registry. Deprecated: use RegisterFunc.

```go
func Register(m Module)
```

### `RegisterFunc`

`RegisterFunc` registers a module constructor that runs once per platform.

```go
func RegisterFunc(f func() Module)
```

### `SetupConnections`

`SetupConnections` will parse the env for named connection strings.

```go
func SetupConnections(environment []string)
```

### `Start`

`Start` is a shorthand to create a new *Platform instance and immediately starts the server listener and handles requests.

```go
func Start(ctx context.Context, options *Options) (*Platform, error)
```

### `TestMiddleware`

`TestMiddleware` returns a middleware that just passes along the request.

```go
func TestMiddleware() Middleware
```

### `Transaction`

`Transaction` wraps a function in a transaction. If the function returns an error, the transaction is rolled back. If the function returns nil, the transaction is committed.

```go
func Transaction(ctx context.Context, db *sql.DB, fn func(context.Context, *sql.Tx) error) error
```

### `URLParam`

`URLParam` will return a named parameter value from the request URL.

```go
func URLParam(r *http.Request, name string) string
```

### `Use`

`Use` adds a middleware to the platform router. Call it in main, not a test.

```go
func Use(mw Middleware)
```

### `Context`

`Context` returns the cancellation context, done once the manager stops.

```go
func (*Manager) Context() context.Context
```

### `Platform`

`Platform` returns the generation currently serving, and nil when there is none: before Start, after Stop, and after a reload that failed.

```go
func (*Manager) Platform() *Platform
```

### `Reload`

`Reload` stops the running platform and starts a new one on the same socket, unless Check refuses it first. Generations never overlap, and the module registered as a value must survive a restart.

```go
func (*Manager) Reload(ctx context.Context) error
```

### `Start`

`Start` binds the listener, starts the first platform generation on the socket, writes Options.PidFile and arms the SIGHUP handler. Cancelling ctx stops the manager.

```go
func (*Manager) Start(ctx context.Context) error
```

### `Stop`

`Stop` shuts down the running platform and closes the socket. A stopped manager does not start again.

```go
func (*Manager) Stop()
```

### `URL`

`URL` gives the e2e endpoint URL for requests. A reload does not change it. Before Start succeeds there is no socket and the URL is empty.

```go
func (*Manager) URL() string
```

### `Wait`

`Wait` will pause until the manager is stopped. A reload does not end it.

```go
func (*Manager) Wait()
```

### `Context`

`Context` returns the cancellation context, done once the server stops.

```go
func (*Platform) Context() context.Context
```

### `Find`

`Find` fills target with the module matching the type.

```go
func (*Platform) Find(target any) bool
```

### `Register`

`Register` adds a module to the platform registry.

```go
func (*Platform) Register(m Module)
```

### `Start`

`Start` starts the server, writes Options.PidFile and prints the registered routes. It stops on a cancelled context, SIGINT or SIGTERM.

```go
func (*Platform) Start(ctx context.Context) error
```

### `Stats`

`Stats` will report how many middlewares and plugins are added to the registry.

```go
func (*Platform) Stats() (int, int)
```

### `Stop`

`Stop` will gracefully shutdown the server and then cancel the server context when done.

Stop is an important part of the lifecycle tests. When closing the registry, each plugins Stop function gets invoked in parallel. This enables the plugin to clear background goroutine event loops, or flush a dirty buffer to storage.

Only after the server has fully shut down does the internal context get cancelled.

```go
func (*Platform) Stop()
```

### `URL`

`URL` gives the e2e endpoint URL for requests.

```go
func (*Platform) URL() string
```

### `Use`

`Use` adds a middleware to the platform registry.

```go
func (*Platform) Use(m Middleware)
```

### `Wait`

`Wait` will pause until the server is shut down.

```go
func (*Platform) Wait()
```

### `Cleanup`

`Cleanup` registers a function to run when the registry closes.

```go
func (*Registry) Cleanup(fn func(context.Context))
```

### `Clone`

`Clone` provides a copy of the registry for use in the platform. Modules registered as constructors are built here, one set per clone.

```go
func (*Registry) Clone() *Registry
```

### `Close`

`Close` runs the cleanups, then clears the modules, middleware and cleanups.

```go
func (*Registry) Close(ctx context.Context)
```

### `Find`

`Find` fills target, a pointer or an interface, with a matching module.

```go
func (*Registry) Find(target any) bool
```

### `Register`

`Register` adds a Module to the registry. Deprecated: use RegisterFunc.

```go
func (*Registry) Register(m Module)
```

### `RegisterFunc`

`RegisterFunc` adds a module constructor that Clone calls per platform.

```go
func (*Registry) RegisterFunc(f func() Module)
```

### `Start`

`Start` starts the modules in order and then mounts them. A failure halts the run, and an error ends it.

```go
func (*Registry) Start(ctx context.Context, mux Router, opts *Options) error
```

### `Stats`

`Stats` returns counts for modules and middlewares in the registry.

```go
func (*Registry) Stats() int
```

### `Use`

`Use` adds a Middleware to the registry.

```go
func (*Registry) Use(f Middleware)
```

### `Middleware`

`Middleware` records requests handled by next.

```go
func (*TelemetryModule) Middleware(next http.Handler) http.Handler
```

### `Mount`

`Mount` registers the debug front end on the platform router.

```go
func (*TelemetryModule) Mount(_ context.Context, r Router) error
```

### `Options`

`Options` returns the options the tracer runs on, after the environment applied.

```go
func (*TelemetryModule) Options() oida.Options
```

### `Tracer`

`Tracer` returns the tracer the module records into.

```go
func (*TelemetryModule) Tracer() *oida.Tracer
```

### `Mount`

`Mount` returns nil (no error).

```go
func (UnimplementedModule) Mount(ctx context.Context, r Router) error
```

### `Name`

`Name` returns an empty string.

```go
func (UnimplementedModule) Name() string
```

### `Start`

`Start` returns nil (no error).

```go
func (UnimplementedModule) Start(ctx context.Context) error
```

### `Stop`

`Stop` returns nil (no error).

```go
func (UnimplementedModule) Stop(ctx context.Context) error
```
