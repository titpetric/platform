package platform_test

import (
	"testing"

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
