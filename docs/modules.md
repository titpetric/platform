# Creating with Modules

## Module Contract

A module must implement:

```go
type Module interface {
    Name() string
    Mount(Router) error
    Start() error
    Stop() error
}
```

The functions are ran in this order; If `Mount` runs, `Start` has completed.

- `Start()` — start background goroutines or services.
- `Mount(Router)` — attach HTTP routes.
- `Stop()` — clean up and stop all background work.

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
