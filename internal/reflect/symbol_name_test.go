package reflect_test

import (
	"testing"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/internal/reflect"
	"github.com/titpetric/platform/module/user/storage"
	"github.com/titpetric/platform/telemetry"
)

var SymbolName = reflect.SymbolName

func TestStartAuto(t *testing.T) {
	input := internal.NewDatabaseProvider()

	t.Log(SymbolName(input))
	t.Log(SymbolName(input.Open))

	storage := &storage.UserStorage{}

	t.Log(SymbolName(storage))
	t.Log(SymbolName(storage.Create))

	t.Log(SymbolName(telemetry.StartAuto))
	t.Log(SymbolName(32))

	// If a string is passed, the string is returned
	t.Log(SymbolName("test.start.auto"))
}
