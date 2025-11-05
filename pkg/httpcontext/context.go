package httpcontext

import (
	"context"
	"net/http"
)

type Value[T any] struct {
	Key any
}

func NewValue[T any](key any) *Value[T] {
	return &Value[T]{Key: key}
}

func (v *Value[T]) Get(r *http.Request) T {
	return v.GetContext(r.Context())
}

func (v *Value[T]) GetContext(ctx context.Context) (res T) {
	if val := ctx.Value(v.Key); val != nil {
		res, _ = val.(T)
	}
	return
}


func (v *Value[T]) Set(r *http.Request, val T) *http.Request {
	*r = *r.WithContext(v.SetContext(r.Context(), val))
	return r
}

func (v *Value[T]) SetContext(ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, v.Key, val)
}
