package internal

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
)

// CountRoutes returns the total number of routes, and the total number of known middlewares.
func CountRoutes(r chi.Routes) (int, int) {
	var (
		routeCount = new(atomic.Int32)
		mwCount    = new(atomic.Int32)
	)

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routeCount.Add(1)
		if len(middlewares) > 0 {
			mwCount.Add(int32(len(middlewares)))
		}
		return nil
	})

	return int(routeCount.Load()), int(mwCount.Load())
}

// PrintRoutes will print the number of routes and middlewares, and the routing table.
func PrintRoutes(r chi.Routes) {
	routes, mws := CountRoutes(r)
	log.Printf("[router] registered %d routes and %d middlewares\n", routes, mws)

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s -> %s\n", method, route, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		return nil
	})
}
