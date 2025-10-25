package internal

import (
	"context"
	"net/http"
)

type ContextValue[T any] struct {
	Key any
}

func NewContextValue[T any](key any) *ContextValue[T] {
	return &ContextValue[T]{Key: key}
}

func (v *ContextValue[T]) Get(r *http.Request) (res T) {
	if val := r.Context().Value(v.Key); val != nil {
		res, _ = val.(T)
	}
	return
}

func (v *ContextValue[T]) Set(r *http.Request, val T) *http.Request {
	ctx := context.WithValue(r.Context(), v.Key, val)
	h := r.WithContext(ctx)
	*r = *h
	return h
}
