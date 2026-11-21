package platform

import (
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// swapRegistry replaces the package registry for one test, and the registration
// here never reaches the state another test observes.
func swapRegistry(t *testing.T) {
	t.Helper()

	registry := global.registry
	global.registry = &Registry{}
	t.Cleanup(func() { global.registry = registry })
}

// TestRegister covers the deprecated package registry: a module registered by
// value is visible to the registry the package holds.
func TestRegister(t *testing.T) {
	swapRegistry(t)

	Register(NewUnimplementedModule("first"))
	modules, _ := global.registry.Stats()
	assert.Equal(t, 1, modules)

	Register(NewUnimplementedModule("second"))
	modules, _ = global.registry.Stats()
	assert.Equal(t, 2, modules)
}

// TestRegisterFunc covers the constructor form, which Clone calls per platform.
func TestRegisterFunc(t *testing.T) {
	swapRegistry(t)

	RegisterFunc(func() Module { return NewUnimplementedModule("built") })
	modules, _ := global.registry.Stats()
	assert.Equal(t, 1, modules)
}

// TestUse covers the package middleware list.
func TestUse(t *testing.T) {
	swapRegistry(t)

	Use(TestMiddleware())
	_, middleware := global.registry.Stats()
	assert.Equal(t, 1, middleware)
}

// TestDatabase covers the package provider, which init fills from the environment.
func TestDatabaseProvider(t *testing.T) {
	var provider DatabaseProvider = Database
	assert.NotNil(t, provider)

	provider = global.db
	_, err := provider.Open(t.Context(), "missing")
	assert.Error(t, err)
}
