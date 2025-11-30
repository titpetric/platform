package platform

import (
	"context"
)

// UnimplementedModule implements the module contract.
// The module can embed the type to skip implementing
// any of the bound functions.
type UnimplementedModule struct {
	NameFn  func() string
	StartFn func(context.Context) error
	StopFn  func(context.Context) error
	MountFn func(context.Context, Router) error
}

// Name returns an empty string.
func (m UnimplementedModule) Name() string {
	if m.NameFn != nil {
		return m.NameFn()
	}
	return ""
}

// Start returns nil (no error).
func (m UnimplementedModule) Start(ctx context.Context) error {
	if m.StartFn != nil {
		return m.StartFn(ctx)
	}
	return nil
}

// Stop returns nil (no error).
func (m UnimplementedModule) Stop(ctx context.Context) error {
	if m.StopFn != nil {
		return m.StopFn(ctx)
	}
	return nil
}

// Mount returns nil (no error).
func (m UnimplementedModule) Mount(ctx context.Context, r Router) error {
	if m.MountFn != nil {
		return m.MountFn(ctx, r)
	}
	return nil
}

var _ Module = UnimplementedModule{}
