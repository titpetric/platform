package internal

import (
	"net/http"
	"testing"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal/assert"
)

func TestRoutesCount(t *testing.T) {
	r := chi.NewRouter()

	routes_a, mws_a := CountRoutes(r)
	assert.Equal(t, 0, routes_a)
	assert.Equal(t, 0, mws_a)

	// middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	// route
	r.Get("/404", http.NotFoundHandler().ServeHTTP)

	routes_b, mws_b := CountRoutes(r)
	assert.Equal(t, 1, routes_b)
	assert.Equal(t, 1, mws_b)

	log := &testLogger{}
	PrintRoutes(log, r)

	// The summary line, and one line for the single route.
	assert.Equal(t, 2, len(log.messages))
	assert.Equal(t, "routes registered", log.messages[0])
	assert.Equal(t, "route", log.messages[1])
}

// testLogger records the messages PrintRoutes writes.
type testLogger struct {
	messages []string
}

func (l *testLogger) Info(msg string, _ ...any) {
	l.messages = append(l.messages, msg)
}
