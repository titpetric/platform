package expvar

import (
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/titpetric/platform/internal/require"
)

func TestHandler(t *testing.T) {
	h := NewHandler()

	require.NotNil(t, h)
	require.NoError(t, h.Mount(chi.NewRouter()))
}
