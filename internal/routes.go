package internal

import (
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
			"handler", runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		return nil
	})
}
