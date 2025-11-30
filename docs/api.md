# Package platform

```go
import (
	"github.com/titpetric/platform"
}
```

The platform is an extensible modular system for writing HTTP servers.

1. Provides a global registry for middleware and module registration
2. Provides a lifecycle to the modules for graceful shutdown
3. Provides a router the modules can attach to

It's advised to use `platform.Register` from `init` functions.
Similarly, `platform.Use` should be used from `main` or any
descendant setup functions. Don't use these functions from tests
as they create a shared state.

It's possible to use the platform in an emperative way.

```go
svc := platform.New(context.Background())
svg.Use(middleware.Logger)
svc.Register(user.NewModule())
```

The platform lifecycle is extensively tested to ensure no races, no
goroutine leaks. Each platform object creates a copy of the global
state and holds scoped allocations only, enabling test parallelism.

## Types

```go
// DatabaseProvider is the implementation interface for working with named connections.
// If no connection name is passed, the "default" connection will be used.
type DatabaseProvider interface {
	Open(...string) (*sqlx.DB, error)
	Connect(context.Context, ...string) (*sqlx.DB, error)
}
```

```go
// ErrorResponse is our JSON choice of an error response.
type ErrorResponse struct {
	Error ErrorResponseBody `json:"error"`
}
```

```go
// ErrorResponseBody is an inner type for ErrorResponse.Error.
type ErrorResponseBody struct {
	Code	int	`json:"code"`
	Message	string	`json:"message"`
}
```

```go
// Middleware is a type alias for middleware functions.
type Middleware func(http.Handler) http.Handler
```

```go
// Module is the implementation contract for modules.
//
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

```go
// Options is a configuration struct for platform behaviour.
type Options struct {
	// ServerAddr is the address the server listens to.
	ServerAddr	string

	// Quiet turns down the verbosity in the Platform logging code, set to true in tests.
	Quiet	bool
}
```

```go
// Platform is our world struct.
type Platform struct {
	options	*Options

	// server setup
	router		*chi.Mux
	server		*http.Server
	listener	net.Listener

	// final shutdown context
	context	context.Context
	cancel	context.CancelFunc
	stop	func()
	once	sync.Once

	// registry holds settings for plugins and middleware.
	// It's auto-filled from global scope.
	registry	*Registry
}
```

```go
// Registry provides a programmatic API to manage middleware and plugins.
// A plugin registers middleware and has a contract to enforce lifecycle.
type Registry struct {
	mu	sync.RWMutex

	modules		[]Module
	middleware	[]Middleware
}
```

```go
// Router is a local shim that aliases the chi router interface.
type Router = chi.Router
```

```go
// UnimplementedModule implements the module contract.
// The module can embed the type to skip implementing
// any of the bound functions.
type UnimplementedModule struct {
	NameFn	func() string
	StartFn	func(context.Context) error
	StopFn	func(context.Context) error
	MountFn	func(context.Context, Router) error
}
```

## Vars

```go
// Database is a holder of the database provider api in package namespace.
var Database DatabaseProvider = global.db
```

## Function symbols

- `func Error (w http.ResponseWriter, r *http.Request, status int, data error)`
- `func FromContext (ctx context.Context) *Platform`
- `func FromRequest (r *http.Request) *Platform`
- `func JSON (w http.ResponseWriter, r *http.Request, status int, data any)`
- `func New (options *Options) *Platform`
- `func NewOptions () *Options`
- `func NewTestOptions () *Options`
- `func Param (r *http.Request, name string) string`
- `func QueryParam (r *http.Request, name string) string`
- `func Register (m Module)`
- `func Start (ctx context.Context, options *Options) (*Platform, error)`
- `func TestMiddleware () Middleware`
- `func Transaction (ctx context.Context, db *sqlx.DB, fn func(context.Context, *sqlx.Tx) error) error`
- `func URLParam (r *http.Request, name string) string`
- `func Use (mw Middleware)`
- `func (*Platform) Context () context.Context`
- `func (*Platform) Find (target any) bool`
- `func (*Platform) Register (m Module)`
- `func (*Platform) Start (ctx context.Context) error`
- `func (*Platform) Stats () (int, int)`
- `func (*Platform) Stop ()`
- `func (*Platform) URL () string`
- `func (*Platform) Use (m Middleware)`
- `func (*Platform) Wait ()`
- `func (*Registry) Clone () *Registry`
- `func (*Registry) Close ()`
- `func (*Registry) Find (target any) bool`
- `func (*Registry) Register (m Module)`
- `func (*Registry) Start (ctx context.Context, mux Router) error`
- `func (*Registry) Stats () int`
- `func (*Registry) Use (f Middleware)`
- `func (UnimplementedModule) Mount (ctx context.Context, r Router) error`
- `func (UnimplementedModule) Name () string`
- `func (UnimplementedModule) Start (ctx context.Context) error`
- `func (UnimplementedModule) Stop (ctx context.Context) error`

### Error

Error writes an error payload as JSON.

```go
func Error (w http.ResponseWriter, r *http.Request, status int, data error)
```

### FromContext

FromContext returns the *Platform instance attached to the context.

```go
func FromContext (ctx context.Context) *Platform
```

### FromRequest

FromRequest returns the *Platform instance attached to the request.

```go
func FromRequest (r *http.Request) *Platform
```

### JSON

JSON writes any payload as JSON. If the payload is nil, the write is omitted.
If an error occurs in encoding, a telemetry error is logged.

```go
func JSON (w http.ResponseWriter, r *http.Request, status int, data any)
```

### New

New will create a new *Platform object. It is the allocation point
for each platform instance. If no options are passed, the defaults are in use.
The defaults options are provided by NewOptions().

```go
func New (options *Options) *Platform
```

### NewOptions

NewOptions provides default options for the platform.

```go
func NewOptions () *Options
```

### NewTestOptions

NewTestOptions produces default options for tests.

```go
func NewTestOptions () *Options
```

### Param

Param will return a named URL parameter, or query string.

```go
func Param (r *http.Request, name string) string
```

### QueryParam

QueryParam will return a named query parameter from the request.

```go
func QueryParam (r *http.Request, name string) string
```

### Register

Register will register a module in the platform global registry.
It should not be relied upon in tests, keeping global state empty.
This enables registering modules using blank imports.

```go
func Register (m Module)
```

### Start

Start is a shorthand to create a new *Platform instance and
immediately starts the server listener and handles requests.

```go
func Start (ctx context.Context, options *Options) (*Platform, error)
```

### TestMiddleware

TestMiddleware returns a middleware that just passes along the request.

```go
func TestMiddleware () Middleware
```

### Transaction

Transaction wraps a function in a transaction.
If the function returns an error, the transaction is rolled back.
If the function returns nil, the transaction is committed.

```go
func Transaction (ctx context.Context, db *sqlx.DB, fn func(context.Context, *sqlx.Tx) error) error
```

### URLParam

URLParam will return a named parameter value from the request URL.

```go
func URLParam (r *http.Request, name string) string
```

### Use

Use will add a middleware to the platform router.
It should not be relied upon in tests, keeping global state empty.
This should be used from main() to define any global middleware.

```go
func Use (mw Middleware)
```

### Context

Context returns the cancellation context for the service.
When the context finishes, the server has shut down.

```go
func (*Platform) Context () context.Context
```

### Find

Find fills target with the module matching the type.

```go
func (*Platform) Find (target any) bool
```

### Register

Register will add a registry.Module into the internal platform registry.
This function should be called before Serve is called.

```go
func (*Platform) Register (m Module)
```

### Start

Start will start the server and print the registered routes.
It respects cancellation from the passed context, as well as
sets up signal notification to respond to SIGTERM.

```go
func (*Platform) Start (ctx context.Context) error
```

### Stats

Stats will report how many middlewares and plugins are added to the registry.

```go
func (*Platform) Stats () (int, int)
```

### Stop

Stop will gracefully shutdown the server and then cancel the server context when done.

Stop is an important part of the lifecycle tests. When closing the registry,
each plugins Stop function gets invoked in parallel. This enables the plugin
to clear background goroutine event loops, or flush a dirty buffer to storage.

Only after the server has fully shut down does the internal context get cancelled.

```go
func (*Platform) Stop ()
```

### URL

URL gives the e2e endpoint URL for requests.

```go
func (*Platform) URL () string
```

### Use

Use will add a middleware to the internal platform registry.
This function should be called before Serve is called.

```go
func (*Platform) Use (m Middleware)
```

### Wait

Wait will pause until the server is shut down.

```go
func (*Platform) Wait ()
```

### Clone

Clone provides a copy of the registry for use in the platform.
Closing the copy leaves the package global state alone.

```go
func (*Registry) Clone () *Registry
```

### Close

Close will invoke all the modules close functions in parallel.
When finished, it will clear the registered modules list, as
well as any defined middleware.

```go
func (*Registry) Close ()
```

### Find

Find gets a Module from the registry.
The target argument can be a pointer or an interface. The function returns true
if a module matching the type or interface was found and assigned to `target`.

```go
func (*Registry) Find (target any) bool
```

### Register

Register adds a Module to the registry.

```go
func (*Registry) Register (m Module)
```

### Start

Start will invoke all the modules start functions sequentially.
If an error occurs, execution is halted and an error is returned.
The context is passed along for observability and access to the platform.

```go
func (*Registry) Start (ctx context.Context, mux Router) error
```

### Stats

Stats returns counts for modules and middlewares in the registry.

```go
func (*Registry) Stats () int
```

### Use

Use adds a Middleware to the registry.

```go
func (*Registry) Use (f Middleware)
```

### Mount

Mount returns nil (no error).

```go
func (UnimplementedModule) Mount (ctx context.Context, r Router) error
```

### Name

Name returns an empty string.

```go
func (UnimplementedModule) Name () string
```

### Start

Start returns nil (no error).

```go
func (UnimplementedModule) Start (ctx context.Context) error
```

### Stop

Stop returns nil (no error).

```go
func (UnimplementedModule) Stop (ctx context.Context) error
```


