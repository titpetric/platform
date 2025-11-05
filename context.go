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
