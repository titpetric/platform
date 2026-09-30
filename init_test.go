package platform

import (
	"testing"

	"github.com/titpetric/platform/internal/assert"
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

	setupConnections(env, collect)

	want := map[string]string{
		"xxx":     "sqlite://:memory:",
		"default": "sqlite://:memory:",
	}

	assert.Equal(t, want, got)
}

// TestDatabaseEnvMalformed covers an environment entry with no "=" or an empty
// name, which SetupConnections accepts from a caller and skips.
func TestDatabaseEnvMalformed(t *testing.T) {
	env := []string{
		"PLATFORM_DB_MAIN",
		"PLATFORM_DB_",
		"PLATFORM_DB_=orphan",
		"PLATFORM_DB_USERS=sqlite://:memory:",
		"UNRELATED=value",
	}

	got := map[string]string{}
	setupConnections(env, func(key, value string) {
		got[key] = value
	})

	want := map[string]string{
		"default": "sqlite://:memory:",
		"users":   "sqlite://:memory:",
	}

	assert.Equal(t, want, got)
}
