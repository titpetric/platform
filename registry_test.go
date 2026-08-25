package platform_test

import (
	"context"
	"strconv"
	"testing"

	chi "github.com/go-chi/chi/v5"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/require"
)

type Greeter interface {
	Greet() string
}

type MyGreeter struct {
	platform.UnimplementedModule

	Nickname string
}

func (g MyGreeter) Greet() string { return "hello, " + g.Nickname }

func TestRegistry_Modules(t *testing.T) {
	r := platform.Registry{}
	r.Register(MyGreeter{Nickname: "Alice"})
	r.Register(&MyGreeter{Nickname: "Bob"})

	t.Run("find into interface (Greeter) sets to first implementing module", func(t *testing.T) {
		var gi Greeter
		ok := r.Find(&gi)
		require.True(t, ok, "Find should succeed for interface target")
		require.NotNil(t, gi, "interface target should be set")
		require.Equal(t, "hello, Alice", gi.Greet())
	})

	t.Run("find into concrete value sets to first matching concrete value", func(t *testing.T) {
		var mg MyGreeter
		ok := r.Find(&mg)
		require.True(t, ok, "Find should succeed for concrete value target")
		require.Equal(t, "Alice", mg.Nickname)
	})

	t.Run("find into concrete pointer picks pointer-registered module", func(t *testing.T) {
		var mpg *MyGreeter // nil
		ok := r.Find(&mpg)
		require.True(t, ok, "Find should succeed for pointer target")
		require.NotNil(t, mpg, "pointer should be set")
		require.Equal(t, "Bob", mpg.Nickname)
	})

	t.Run("non-pointer target returns false", func(t *testing.T) {
		var notPtr MyGreeter
		ok := r.Find(notPtr)
		require.False(t, ok, "Find should fail for non-pointer target")
	})

	t.Run("nil interface target returns false", func(t *testing.T) {
		require.False(t, r.Find(nil), "Find(nil) should return false")
	})
}

// TestRegistry_RegisterFunc covers constructor registration, which is what
// gives each platform, and each reload generation, modules of its own.
func TestRegistry_RegisterFunc(t *testing.T) {
	var built int

	r := platform.Registry{}
	r.RegisterFunc(func() platform.Module {
		built++
		return &MyGreeter{Nickname: strconv.Itoa(built)}
	})

	modules, _ := r.Stats()
	require.Equal(t, 1, modules)
	require.Equal(t, 0, built, "a constructor is called by Clone, not by registration")

	var first, second *MyGreeter
	require.True(t, r.Clone().Find(&first))
	require.True(t, r.Clone().Find(&second))

	require.Equal(t, 2, built)
	require.True(t, first != second, "each clone should get its own module")
	require.Equal(t, "1", first.Nickname)
	require.Equal(t, "2", second.Nickname)
}

// TestRegistry_Register covers the deprecated value form, where the module
// is shared by everything that clones the registry.
func TestRegistry_Register(t *testing.T) {
	value := &MyGreeter{Nickname: "shared"}

	r := platform.Registry{}
	r.Register(value)

	var first, second *MyGreeter
	require.True(t, r.Clone().Find(&first))
	require.True(t, r.Clone().Find(&second))

	require.True(t, first == value)
	require.True(t, second == value)
}

// TestRegistry_start_materializes covers a registry started as it stands,
// with no clone to call the constructors it holds.
func TestRegistry_start_materializes(t *testing.T) {
	var started int

	r := platform.Registry{}
	r.RegisterFunc(func() platform.Module {
		mod := platform.NewUnimplementedModule("materialized")
		mod.StartFn = func(context.Context) error {
			started++
			return nil
		}
		return mod
	})

	require.NoError(t, r.Start(t.Context(), chi.NewRouter(), platform.NewTestOptions()))
	require.Equal(t, 1, started)
}
