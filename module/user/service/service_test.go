package service

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestService_Mount(t *testing.T) {
	s := &Service{}
	s.Mount(chi.NewRouter())
}
