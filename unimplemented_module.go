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
	StopFn  func() error
	MountFn func(Router) error
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
func (m UnimplementedModule) Stop() error {
	if m.StopFn != nil {
		return m.StopFn()
	}
	return nil
}

// Mount returns nil (no error).
func (m UnimplementedModule) Mount(r Router) error {
	if m.MountFn != nil {
		return m.MountFn(r)
	}
	return nil
}

var _ Module = UnimplementedModule{}
