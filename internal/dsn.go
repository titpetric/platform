package internal

import (
	"net/url"
	"strings"
)

func cleanDSN(driver, dsn string) string {
	switch driver {
	case "sqlite":
		if isSQLiteMemoryDSN(dsn) {
			return dsn
		}

		head, query := splitDSN(dsn)
		query = addOption(query, "_busy_timeout", "5000")
		query = addOption(query, "_journal_mode", "wal")
		return joinDSN(head, query)
	case "mysql":
		head, query := splitDSN(dsn)
		query = addOption(query, "collation", "utf8mb4_general_ci")
		query = addOption(query, "parseTime", "true")
		query = addOption(query, "loc", "Local")
		return joinDSN(head, query)
	default:
		return dsn
	}
}

// splitDSN separates a DSN into the part before its query and the query
// itself, neither carrying the "?".
//
// The query begins at the first "?" at or after the last "/". Every DSN form
// the platform accepts puts the query last and cannot have a "/" in it: the
// mysql grammar is [user[:pass]@][net[(addr)]]/dbname[?params], and a sqlite
// DSN is a path or a file: URI. So a "?" earlier than that belongs to a
// password or to a directory name, and is not a query at all. Scanning the
// whole string for "?" instead turns user:pa?ss@tcp(h)/db into a connection to
// a database named db&collation=... , because the options are then appended
// with "&" to something that has no query to extend.
func splitDSN(dsn string) (head, query string) {
	from := strings.LastIndex(dsn, "/") + 1

	if i := strings.Index(dsn[from:], "?"); i >= 0 {
		i += from
		return dsn[:i], dsn[i+1:]
	}

	return dsn, ""
}

// joinDSN is the inverse of splitDSN.
func joinDSN(head, query string) string {
	if query == "" {
		return head
	}
	return head + "?" + query
}

// addOption appends key=value unless the query already sets key. The existing
// query is never re-encoded: a driver that accepts bytes url would escape
// differently keeps the DSN it was given.
func addOption(query, key, value string) string {
	if hasOption(query, key) {
		return query
	}
	return joinOption(query, key+"="+value)
}

func joinOption(query, option string) string {
	if query == "" {
		return option
	}
	return query + "&" + option
}

// hasOption reports whether the query already sets key. A query the standard
// library cannot parse falls back to a textual match, so a DSN a driver
// accepts but url does not is still never given a duplicate option.
func hasOption(query, key string) bool {
	if values, err := url.ParseQuery(query); err == nil {
		return values.Has(key)
	}
	return strings.Contains(query, key+"=")
}

func isSQLiteMemoryDSN(dsn string) bool {
	database, query := splitDSN(dsn)
	if strings.EqualFold(database, ":memory:") || strings.EqualFold(database, "file::memory:") {
		return true
	}

	values, _ := url.ParseQuery(query)
	return strings.EqualFold(values.Get("mode"), "memory")
}
