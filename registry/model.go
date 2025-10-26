package registry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router is a local shim that aliases the chi router interface.
// This is also aliased to `platform.Router`, to not rely on importing this package.
type Router = chi.Router

// MiddlewareFunc is a type alias for middlewares.
type MiddlewareFunc func(http.Handler) http.Handler
