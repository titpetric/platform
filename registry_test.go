package platform

import (
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestRegistryUse covers the middleware list a registry holds.
func TestRegistryUse(t *testing.T) {
	var registry *Registry = &Registry{}

	registry.Use(TestMiddleware())
	modules, middleware := registry.Stats()
	assert.Equal(t, 0, modules)
	assert.Equal(t, 1, middleware)
}

// TestRegistryRegister covers a module registered by value.
func TestRegistryRegister(t *testing.T) {
	registry := &Registry{}
	registry.Register(NewUnimplementedModule("by-value"))

	modules, _ := registry.Stats()
	assert.Equal(t, 1, modules)
}

// TestRegistryRegisterFunc covers a module registered by constructor, which
// Clone calls for the platform it builds.
func TestRegistryRegisterFunc(t *testing.T) {
	registry := &Registry{}
	registry.RegisterFunc(func() Module { return NewUnimplementedModule("constructed") })

	clone := registry.Clone()
	modules, _ := clone.Stats()
	assert.Equal(t, 1, modules)

	var found Module
	assert.True(t, clone.Find(&found))
	assert.Equal(t, "constructed", found.Name())
}
