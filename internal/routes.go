package internal

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	chi "github.com/go-chi/chi/v5"
)

// Logger is the subset of *slog.Logger that PrintRoutes writes through.
// The platform logger satisfies it, and internal can't import the platform
// package to name it.
type Logger interface {
	Info(msg string, args ...any)
}

// CountRoutes returns the total number of routes, and the total number of known middlewares.
func CountRoutes(r chi.Routes) (int, int) {
	var mu sync.Mutex
	var routes, mws int

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		mu.Lock()
		defer mu.Unlock()

		routes++
		if len(middlewares) > 0 {
			mws += len(middlewares)
		}
		return nil
	})

	return routes, mws
}

// PrintRoutes will print the number of routes and middlewares, and the routing table.
func PrintRoutes(log Logger, r chi.Routes) {
	routes, mws := CountRoutes(r)
	log.Info("routes registered", "routes", routes, "middlewares", mws)

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Info("route", "method", method, "path", route,
			"handler", handlerName(handler))
		return nil
	})
}

// handlerName names a handler for the route log, resolving the symbol where
// the handler has one and falling back to its type. reflect.Value.Pointer is
// legal only for a chan, func, map, pointer, slice or unsafe pointer, and a
// handler may be a struct value with a ServeHTTP method, so the kind decides.
func handlerName(handler http.Handler) string {
	value := reflect.ValueOf(handler)

	switch value.Kind() {
	case reflect.Func:
		if fn := runtime.FuncForPC(value.Pointer()); fn != nil {
			return fn.Name()
		}
	case reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		if value.IsNil() {
			break
		}
		if fn := runtime.FuncForPC(value.Pointer()); fn != nil {
			return fn.Name()
		}
	}

	return fmt.Sprintf("%T", handler)
}
