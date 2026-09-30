# Creating Modules

## Module Contract

A module must implement:

```go
type Module interface {
	Name() string
	Start(context.Context) error
	Mount(context.Context, Router) error
	Stop(context.Context) error
}
```

The functions run in this order, and each phase completes across every module before the next begins: every `Start` has returned before the first `Mount` runs, so a module may rely in `Mount` on state another module built in `Start`. Within a phase the order is registration order.

- `Start(context.Context) error` - start background goroutines or services.
- `Mount(context.Context, Router) error` - attach HTTP routes. Middleware belongs in `platform.Use` or `(*Platform).Use`, not here: `r.Use` from `Mount` panics once another module has registered a route.
- `Stop(context.Context) error` - clean up and stop all background work. Every module's `Stop` runs in parallel on shutdown, so it must not depend on another module's teardown, and it runs for a module whose own `Start` returned an error, so it has to tolerate a partially built value. A panic in `Stop` is recovered and recorded.

The context `Stop` receives is the platform's shutdown context, already cancelled by the time it arrives, so a `Stop` that needs a deadline has to make its own.

### Firewalling modules

The context passed to the `Start` function may be used to invoke `platform.FromContext(ctx) *Platform`. This may then use the `Find` api to get a reference to a side-loaded module. It can be used with interfaces.

An example of that would be a user module that provides a certain API.

```go
type SessionService interface {
	IsLoggedIn(context.Context) bool
	GetSessionUser(context.Context) (*model.User, error)
}

func (m *Module) Start(ctx context.Context) error {
	var api SessionService
	if !platform.FromContext(ctx).Find(&api) {
		return errors.New("session service not registered")
	}
	m.session = api
	return nil
}
```

`Find` resolves against registrations, not against started modules, so it can hand back a module whose `Start` has not run yet. A module that needs its dependency ready should use it from `Mount` or from a request rather than from `Start`.

Or it may expose its complete storage API:

```go
type UserService interface {
	SessionStorage() model.SessionStorage
	UserStorage() model.UserStorage
}
```

This allows API usage behind interfaces. In our case, we can expose module-scoped functionality from the individual modules. It's also possible to get a concrete `*user.Handler` type (no firewall).

### Using `UnimplementedModule`

Embed `platform.UnimplementedModule` to reduce boilerplate and override only the methods you need. `Name` is the exception: it has to return a non-empty string or startup fails with `module %T doesn't return name`. Override it, or embed the pointer `NewUnimplementedModule` returns, which fills `NameFn`:

```go
type StaticModule struct {
	platform.UnimplementedModule
}

func (m *StaticModule) Name() string { return "static" }

func (m *StaticModule) Mount(_ context.Context, r platform.Router) error {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("static page"))
	})
	return nil
}
```

The type also carries a hook per method, `NameFn`, `StartFn`, `StopFn` and `MountFn`. Assigning those is the shorter route for a test or a module small enough not to want a type of its own.

## Minimal App Example

```go
func main() {
	// Register common middleware.
	platform.Use(loggingMiddleware)
	platform.RegisterFunc(func() platform.Module { return &StaticModule{} })

	p, err := platform.Start(context.Background(), platform.NewOptions())
	if err != nil {
		log.Fatalf("exit error: %v", err)
	}

	p.Wait()
}
```

`platform.Start` returns as soon as the server is accepting, so `Wait` is what blocks until shutdown. To get `SIGHUP` reloads instead, run a `platform.Manager` or call `cmd.Main`; see [The Platform](platform.md).

`loggingMiddleware` example:

```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		platform.FromRequest(r).Logger.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```
