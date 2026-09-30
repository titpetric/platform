package platform

import (
	"os"
	"strings"
)

func init() {
	SetupConnections(os.Environ())
}

// SetupConnections will parse the env for named connection strings.
func SetupConnections(environment []string) {
	setupConnections(environment, global.db.Register)
}

func setupConnections(environment []string, register func(string, string)) {
	connections := map[string]string{
		"default": "sqlite://:memory:",
	}

	for _, e := range environment {
		if clean, ok := strings.CutPrefix(e, "PLATFORM_DB_"); ok {
			name, dsn, ok := strings.Cut(clean, "=")
			if !ok || name == "" {
				// An entry with no "=" is not a variable. os.Environ never
				// produces one, but SetupConnections is exported and an
				// environment assembled for an exec'd process can.
				continue
			}

			connections[strings.ToLower(name)] = dsn
		}
	}

	for name, dsn := range connections {
		register(name, dsn)
	}
}
