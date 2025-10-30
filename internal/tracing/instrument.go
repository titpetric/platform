// The tracing package will wrap the sqlite driver in use via `init` with elastic apm instrumentation.
// We can't use the upstream module since they implement a different sqlite3 driver provider.
package tracing

import (
	"modernc.org/sqlite"

	"github.com/go-sql-driver/mysql"
	"go.elastic.co/apm/module/apmsql/v2"
)

func init() {
	apmsql.Register("sqlite", &sqlite.Driver{}, apmsql.WithDSNParser(ParseDSN))
	apmsql.Register("mysql", &mysql.MySQLDriver{}, apmsql.WithDSNParser(ParseDSN))
}
