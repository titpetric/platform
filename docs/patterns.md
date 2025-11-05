# Common Patterns

## Mounting Routes

Attach GET and POST endpoints:

```go
func (m *Module) Mount(r platform.Router) error {
    r.Get("/items", m.GetItems)
    r.Post("/items", m.PostItem)
    return nil
}
```

The `platform.Router` is an alias of `chi.Router` (v5). It allows you to
use any of the methods defined in the interface.

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

## Background jobs

The module can implement it's background job lifecycle by providing a
`Start` and `Stop` function. Invoking `Stop` should be a blocking
operation. For example, with `robfig/cron`:

```go
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

func (c *Crontab) Stop() error {
	<-c.scheduler.Stop().Done()
	return nil
}
```

Since `Stop` is blocking, it will wait up to 3 seconds here, so that any
running scheduled task is completed before exiting.


## Middleware

- Add global middleware via `platform.Use()` (package) or `(*Platform).Use()` (instance).
- Middleware should be added **before** `Start(context.Context)`.
- You can use any existing middleware as long as it implements the `Middleware` interface.
