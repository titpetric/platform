package platform

import (
	"os"
	"strings"
)

var Database *DatabaseFactory = &DatabaseFactory{
	credentials: make(map[string]string),
}

func init() {
	setupConnections(Database.Add)
}

func setupConnections(add func(string, string)) {
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
		add(name, dsn)
	}
}
