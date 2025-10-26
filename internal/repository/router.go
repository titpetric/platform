package repository

import (
	"github.com/go-chi/chi/v5"
)

// Router is a local shim that aliases the chi router interface.
type Router = chi.Router
