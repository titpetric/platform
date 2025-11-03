package reflect_test

import (
	"testing"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/pkg/reflect"
	"github.com/titpetric/platform/pkg/require"
	"github.com/titpetric/platform/pkg/telemetry"
)

var SymbolName = reflect.SymbolName

func TestStartAuto(t *testing.T) {
	input := internal.NewDatabaseProvider(nil)

	// assert expected symbol location
	require.Equal(t, "internal.DatabaseProvider", SymbolName(input))
	require.Equal(t, "internal.DatabaseProvider.Open", SymbolName(input.Open))

	// cross package scope doesn't change the underlying type
	require.Equal(t, "internal.DatabaseProvider", SymbolName(platform.Database))

	// interface reference changes the path as defined in the interface
	require.Equal(t, "platform.DatabaseProvider.Connect", SymbolName(platform.Database.Connect))

	// global functions work
	require.Equal(t, "telemetry.StartAuto", SymbolName(telemetry.StartAuto))

	// native types
	require.Equal(t, "int", SymbolName(32))

	// a string will be returned as is
	require.Equal(t, "test.start.auto", SymbolName("test.start.auto"))

	// if a string has slashes, it will be trimmed to remove everything until the last `/`.
	require.Equal(t, "test.start.auto", SymbolName("github.com/titpetric/internal/test.start.auto"))

	// nil is kind of an illegal value (untyped) but here we are
	require.Equal(t, "<nil>", SymbolName(nil))
}
