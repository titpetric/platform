package platform

import (
	"net/http"
)

// Middleware is a type alias for middleware functions.
type Middleware func(http.Handler) http.Handler

// TestMiddleware returns a middleware that just passes along the request.
func TestMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
