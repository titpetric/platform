package internal

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	chi "github.com/go-chi/chi/v5"
)

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
func PrintRoutes(r chi.Routes) {
	routes, mws := CountRoutes(r)
	log.Printf("[router] registered %d routes and %d middlewares\n", routes, mws)

	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s -> %s\n", method, route, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		return nil
	})
}
