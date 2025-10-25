package internal

import (
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/go-chi/chi/v5"
)

func PrintRoutes(r chi.Routes) {
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s -> %s\n", method, route, runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		return nil
	})
}
