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

// splitDSN separates a DSN into the part before its query and the query.
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

// addOption appends key=value unless the query already sets key.
func addOption(query, key, value string) string {
	if hasOption(query, key) {
		return query
	}
	if query == "" {
		return key + "=" + value
	}
	return query + "&" + key + "=" + value
}

// hasOption reports whether the query sets key, by parse or by text.
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
