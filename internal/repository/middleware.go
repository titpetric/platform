package repository

import (
	"net/http"
)

// Middleware is a type alias for middleware functions.
type Middleware func(http.Handler) http.Handler
