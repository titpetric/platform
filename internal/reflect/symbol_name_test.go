package reflect_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/reflect"
	"github.com/titpetric/platform/module/user/storage"
	"github.com/titpetric/platform/telemetry"
)

var SymbolName = reflect.SymbolName

func TestStartAuto(t *testing.T) {
	input := internal.NewDatabaseProvider()

	require.Equal(t, "internal.DatabaseProvider", SymbolName(input))
	require.Equal(t, "internal.DatabaseProvider.Open", SymbolName(input.Open))

	storage := &storage.UserStorage{}

	require.Equal(t, "storage.UserStorage", SymbolName(storage))
	require.Equal(t, "storage.UserStorage.Create", SymbolName(storage.Create))

	require.Equal(t, "telemetry.StartAuto", SymbolName(telemetry.StartAuto))
	require.Equal(t, "int", SymbolName(32))

	require.Equal(t, "test.start.auto", SymbolName("test.start.auto"))
	require.Equal(t, "test.start.auto", SymbolName("github.com/titpetric/internal/test.start.auto"))
	require.Equal(t, "<nil>", SymbolName(nil))
}
