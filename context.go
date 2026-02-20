package platform

import (
	"context"
	"net/http"

	"github.com/titpetric/platform/pkg/httpcontext"
)

type platformKey struct{}

var platformContext = httpcontext.NewValue[*Platform](platformKey{})

// FromRequest returns the *Platform instance attached to the request.
func FromRequest(r *http.Request) *Platform {
	return platformContext.Get(r)
}

// FromContext returns the *Platform instance attached to the context.
func FromContext(ctx context.Context) *Platform {
	return platformContext.GetContext(ctx)
}

type optionsKey struct{}

var optionsContext = httpcontext.NewValue[*Options](optionsKey{})

// OptionsFromRequest returns the *Options instance attached to the request.
func OptionsFromRequest(r *http.Request) *Options {
	return optionsContext.Get(r)
}

// OptionsFromContext returns the *Options instance attached to the context.
func OptionsFromContext(ctx context.Context) *Options {
	return optionsContext.GetContext(ctx)
}
