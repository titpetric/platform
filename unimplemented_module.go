package platform

// UnimplementedModule implements the module contract.
// The module can embed the type to skip implementing
// any of the bound functions.
type UnimplementedModule struct{}

// Name returns an empty string.
func (UnimplementedModule) Name() string {
	return ""
}

// Start returns nil (no error).
func (UnimplementedModule) Start() error {
	return nil
}

// Stop returns nil (no error).
func (UnimplementedModule) Stop() error {
	return nil
}

// Mount returns nil (no error).
func (UnimplementedModule) Mount(Router) error {
	return nil
}

var _ Module = UnimplementedModule{}
