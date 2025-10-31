package service

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/telemetry"
)

// errorMessageKey is a request context scoped value. If an error
// occurs in let's say POST /login, the intent is to set the
// error to the request context, and then render a view to display.
type errorMessageKey struct{}

var errorMessageContext = internal.NewContextValue[string](errorMessageKey{})

func (h *Service) Error(r *http.Request, message string, err error) {
	errorMessageContext.Set(r, message)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if err == nil {
		err = errors.New(message)
	}
	telemetry.CaptureError(r.Context(), err)
}

func (h *Service) GetError(r *http.Request) string {
	return errorMessageContext.Get(r)
}
