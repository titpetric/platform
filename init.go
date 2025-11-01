package platform

import (
	"os"
	"strings"
)

func init() {
	setupConnections(global.db.Register)
}

// setupConnections will parse the env for named connection strings.
func setupConnections(register func(string, string)) {
	connections := map[string]string{
		"default": "sqlite://:memory:",
	}

	for _, e := range os.Environ() {
		if clean, ok := strings.CutPrefix(e, "PLATFORM_DB_"); ok {
			pair := strings.SplitN(clean, "=", 2)

			connections[strings.ToLower(pair[0])] = pair[1]
		}
	}

	for name, dsn := range connections {
		register(name, dsn)
	}
}
