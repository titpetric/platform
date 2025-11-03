package platform

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
)

// TestDatabaseEnv checks that we properly decode the expected environment
// and collect it for a named sql connection string map.
func TestDatabaseEnv(t *testing.T) {
	env := []string{
		"PLATFORM_DB_XXX=sqlite://:memory:",
		"PLATFORM_DB_DEFAULT=sqlite://:memory:",
	}

	got := map[string]string{}
	collect := func(key, value string) {
		got[key] = value
	}

	setupConnections(collect, env)

	want := map[string]string{
		"xxx":     "sqlite://:memory:",
		"default": "sqlite://:memory:",
	}

	require.Equal(t, want, got)
}
