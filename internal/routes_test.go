package internal

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal/require"
)

func TestRoutesCount(t *testing.T) {
	r := chi.NewRouter()

	routes_a, mws_a := CountRoutes(r)
	require.Equal(t, 0, routes_a)
	require.Equal(t, 0, mws_a)

	// middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	// route
	r.Get("/404", http.NotFoundHandler().ServeHTTP)

	routes_b, mws_b := CountRoutes(r)
	require.Equal(t, 1, routes_b)
	require.Equal(t, 1, mws_b)

	PrintRoutes(r)
}
