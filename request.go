package platform

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

// URLParam will return a named parameter value from the request URL.
func URLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

// QueryParam will return a named query parameter from the request.
func QueryParam(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

// Param will return a named URL parameter, or query string.
func Param(r *http.Request, name string) string {
	if val := URLParam(r, name); val != "" {
		return val
	}
	return QueryParam(r, name)
}
