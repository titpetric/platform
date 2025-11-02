package internal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewServer(router chi.Router) *http.Server {
	return &http.Server{
		Handler: router,
	}
}
