package platform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseEnv(t *testing.T) {
	t.Setenv("PLATFORM_DB_XXX", "sqlite://:memory:")
	t.Setenv("PLATFORM_DB_DEFAULT", "sqlite://:memory:")

	got := map[string]string{}
	collect := func(key, value string) {
		got[key] = value
	}

	setupConnections(collect)

	want := map[string]string{
		"xxx":     "sqlite://:memory:",
		"default": "sqlite://:memory:",
	}

	require.Equal(t, want, got)
}
