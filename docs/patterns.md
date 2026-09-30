# Common Patterns

## Mounting Routes

Attach GET and POST endpoints:

```go
func (m *Module) Mount(_ context.Context, r platform.Router) error {
	r.Get("/items", m.GetItems)
	r.Post("/items", m.PostItem)
	return nil
}
```

The `platform.Router` is an alias of `chi.Router` (v5). It allows you to use any of the methods defined in the interface.

Handlers are methods on the module:

```go
func (m *Module) GetItems(w http.ResponseWriter, r *http.Request) {
	// fetch and return items
}

func (m *Module) PostItem(w http.ResponseWriter, r *http.Request) {
	// validate input and create item
}
```

## POST/GET Validation Pattern

For simple validation, parse POST data and call GET handler on error:

```go
func (m *Module) PostItem(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("name") == "" {
		// reuse GET handler to re-render with error
		m.GetItems(w, r)
		return
	}
	// continue processing
}
```

The GET handler has to read the submitted values with `r.FormValue` or `r.PostFormValue`. `platform.Param` and `platform.QueryParam` read the URL path and the query string only, so a handler that uses them re-renders an empty form.

## Background jobs

The module can implement its background job lifecycle by providing a `Start` and `Stop` function. Invoking `Stop` should be a blocking operation. A module still has to satisfy the whole contract, so embed `platform.UnimplementedModule` and override only these two. For example, with `github.com/robfig/cron/v3`:

```go
type Crontab struct {
	*platform.UnimplementedModule
	scheduler *cron.Cron
}

func NewCrontab() platform.Module {
	return &Crontab{
		UnimplementedModule: platform.NewUnimplementedModule("crontab"),
		scheduler:           cron.New(),
	}
}

func (c *Crontab) Start(context.Context) error {
	_, err := c.scheduler.AddFunc("@every 5s", func() {
		log.Printf("This is your cron job starting.")
		time.Sleep(3 * time.Second)
		log.Printf("Cron job exiting after 3 secs.")
	})
	if err != nil {
		return err
	}

	c.scheduler.Start()
	return nil
}

func (c *Crontab) Stop(context.Context) error {
	<-c.scheduler.Stop().Done()
	return nil
}
```

Since `Stop` is blocking, it will wait up to 3 seconds here, so that any running scheduled task is completed before exiting. `Platform.Stop` waits for every module's `Stop` to return, however long that takes. The context it passes carries a five second budget, so a module with nothing better to go on can bound its own teardown by it.

## Middleware

- Add global middleware via `platform.Use()` (package) or `(*Platform).Use()` (instance).
- Package-level `platform.Use()` must be called before `platform.New()` or `platform.Start()`. `New` clones the global registry once, so a later call never reaches that platform. Call it from `main` or an `init`.
- Instance `(*Platform).Use()` must be called before `(*Platform).Start()`. Under a `platform.Manager`, register it from `Manager.Setup`, which runs against every generation.
- Adding middleware too late is silent: no error, no panic, the middleware just never runs.
- You can use any existing middleware as long as it matches the `Middleware` signature, `func(http.Handler) http.Handler`.
