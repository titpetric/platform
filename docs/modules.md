# Creating with Modules

## Module Contract

A module must implement:

```go
type Module interface {
    Name() string
    Mount(Router) error
    Start(context.Context) error
    Stop() error
}
```

The functions are ran in this order; If `Mount` runs, `Start` has completed.

- `Start(context.Context)` — start background goroutines or services.
- `Mount(Router)` — attach HTTP routes.
- `Stop()` — clean up and stop all background work.

### Firewalling modules

The context passed to the `Start` function may be used to invoke
`platform.FromContext(ctx) *Platform`. This may then use the `Find` api to get a
reference to a side-loaded module. It can be used with interfaces.

An example of that would be a user module that provides a certain API.

```go
type SessionService interface {
	IsLoggedIn(context.Context) bool
	GetSessionUser(context.Context) (*model.User, error)
}
var api SessionService
ok := platform.FromContext(r.Context()).Find(&api)
```

Or it may expose it's complete storage API:

```go
type UserService interface {
	SessionStorage() model.SessionStorage
	UserStorage() model.UserStorage
}
```

This allows API usage behind interfaces. In our case, we can expose
module-scoped functionality from the individual modules. It's also
possible to get a concrete `*user.Handler` type (no firewall).

### Using `UnimplementedModule`

Embed `platform.UnimplementedModule` to reduce boilerplate and override only methods you need:

```go
type StaticModule struct {
    platform.UnimplementedModule
}

func (m *StaticModule) Name() string { return "static" }

func (m *StaticModule) Mount(r platform.Router) error {
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("static page"))
    })
    return nil
}
```

## Minimal App Example

```go
func main() {
	// Register common middleware.
	platform.Use(loggingMiddleare)
	platform.Register(&StaticModule{})

	if err := platform.Start(); err != nil {
		log.Fatalf("exit error: %v", err)
	}
}
```

`loggingMiddleware` example:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```
