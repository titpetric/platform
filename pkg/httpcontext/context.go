package httpcontext

import (
	"context"
	"net/http"
)

// Value is used to lookup values stored in context.Context.
type Value[T any] struct {
	Key any
}

// NewValue creates a new value.
func NewValue[T any](key any) *Value[T] {
	return &Value[T]{Key: key}
}

// Get returns a value bound to a request context.
func (v *Value[T]) Get(r *http.Request) T {
	return v.GetContext(r.Context())
}

// GetContext returns a value bound to a context.
func (v *Value[T]) GetContext(ctx context.Context) (res T) {
	if val := ctx.Value(v.Key); val != nil {
		res, _ = val.(T)
	}
	return
}

// Set will update a context value in the passed request.
func (v *Value[T]) Set(r *http.Request, val T) {
	*r = *r.WithContext(v.SetContext(r.Context(), val))
}

func (v *Value[T]) SetContext(ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, v.Key, val)
}
