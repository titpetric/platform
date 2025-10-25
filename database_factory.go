package platform

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DatabaseFactory struct {
	credentials map[string]string
}

func (r *DatabaseFactory) Add(name string, config string) {
	r.credentials[name] = config
}

func (r *DatabaseFactory) Connect(dbName ...string) (*sqlx.DB, error) {
	names := dbName
	if len(names) == 0 {
		names = []string{"default"}
	}

	for _, name := range names {
		if value, ok := r.credentials[name]; ok {
			driver, dsn := "mysql", value

			// allow specifying the driver with url notation,
			// in the follwing form: <driver>://<dsn>.
			if sepIndex := strings.Index(dsn, "://"); sepIndex != -1 {
				driver = dsn[:sepIndex]
				dsn = dsn[sepIndex+3:]
			}

			// we should parse/add these options to the DSN as mandatory:
			// ?collation=utf8_general_ci&parseTime=true&loc=Local
			dsn += "?parseTime=true"
			return sqlx.Connect(driver, dsn)
		}
	}
	return nil, fmt.Errorf("No configuration found for database: %v", names)
}
