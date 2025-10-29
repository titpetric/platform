package service

import (
	"net/http"

	"github.com/titpetric/platform/internal"
)

// errorMessageKey is a request context scoped value. If an error
// occurs in let's say POST /login, the intent is to set the
// error to the request context, and then render a view to display.
type errorMessageKey struct{}

var errorMessageContext = internal.NewContextValue[string](errorMessageKey{})

func (h *Service) Error(r *http.Request, message string, err error) {
	errorMessageContext.Set(r, message)
}

func (h *Service) GetError(r *http.Request) string {
	return errorMessageContext.Get(r)
}
