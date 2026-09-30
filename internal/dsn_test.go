package internal

import (
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

func TestCleanDSN(t *testing.T) {
	type testCase struct {
		name string
		dsn  string
		want string
	}

	tests := []testCase{
		{
			name: "empty DSN",
			dsn:  "",
			want: "?collation=utf8mb4_general_ci&parseTime=true&loc=Local",
		},
		{
			name: "dsn with question mark",
			dsn:  "user:pass@tcp(localhost:3306)/dbname?",
			want: "user:pass@tcp(localhost:3306)/dbname?collation=utf8mb4_general_ci&parseTime=true&loc=Local",
		},
		{
			name: "dsn with collation set",
			dsn:  "user:pass@tcp(localhost:3306)/dbname?collation=utf8",
			want: "user:pass@tcp(localhost:3306)/dbname?collation=utf8&parseTime=true&loc=Local",
		},
		{
			name: "dsn with parseTime set",
			dsn:  "user:pass@tcp(localhost:3306)/dbname?parseTime=false",
			want: "user:pass@tcp(localhost:3306)/dbname?parseTime=false&collation=utf8mb4_general_ci&loc=Local",
		},
		{
			name: "dsn with loc set",
			dsn:  "user:pass@tcp(localhost:3306)/dbname?loc=UTC",
			want: "user:pass@tcp(localhost:3306)/dbname?loc=UTC&collation=utf8mb4_general_ci&parseTime=true",
		},
		{
			name: "dsn with all options set",
			dsn:  "user:pass@tcp(localhost:3306)/dbname?collation=abc&parseTime=abc&loc=abc",
			want: "user:pass@tcp(localhost:3306)/dbname?collation=abc&parseTime=abc&loc=abc",
		},
		{
			// A "?" in the password is not the start of the query. Scanning
			// the whole DSN for one appended the options with "&", which the
			// driver read as part of the database name.
			name: "question mark in the password",
			dsn:  "user:pa?ss@tcp(localhost:3306)/dbname",
			want: "user:pa?ss@tcp(localhost:3306)/dbname?collation=utf8mb4_general_ci&parseTime=true&loc=Local",
		},
		{
			name: "question mark in the password and a real query",
			dsn:  "user:pa?ss@tcp(localhost:3306)/dbname?loc=UTC",
			want: "user:pa?ss@tcp(localhost:3306)/dbname?loc=UTC&collation=utf8mb4_general_ci&parseTime=true",
		},
		{
			name: "question mark in the user",
			dsn:  "us?er:pass@tcp(localhost:3306)/dbname",
			want: "us?er:pass@tcp(localhost:3306)/dbname?collation=utf8mb4_general_ci&parseTime=true&loc=Local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanDSN("mysql", tt.dsn)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCleanDSNDriverSpecific(t *testing.T) {
	dsn := "user=postgres password=secret host=127.0.0.1 port=15432 dbname=postgres sslmode=verify-ca"

	for _, driver := range []string{"postgres", "pgx"} {
		t.Run(driver, func(t *testing.T) {
			assert.Equal(t, dsn, cleanDSN(driver, dsn))
		})
	}
}

func TestCleanSQLiteDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "file database",
			dsn:  "app.db",
			want: "app.db?_busy_timeout=5000&_journal_mode=wal",
		},
		{
			name: "existing query",
			dsn:  "file:app.db?cache=shared",
			want: "file:app.db?cache=shared&_busy_timeout=5000&_journal_mode=wal",
		},
		{
			name: "explicit options",
			dsn:  "app.db?_busy_timeout=1000&_journal_mode=delete",
			want: "app.db?_busy_timeout=1000&_journal_mode=delete",
		},
		{
			name: "missing journal mode",
			dsn:  "app.db?_busy_timeout=1000",
			want: "app.db?_busy_timeout=1000&_journal_mode=wal",
		},
		{
			name: "memory database",
			dsn:  ":memory:",
			want: ":memory:",
		},
		{
			name: "memory URI",
			dsn:  "file::memory:?cache=shared",
			want: "file::memory:?cache=shared",
		},
		{
			name: "named shared memory database",
			dsn:  "file:app?mode=memory&cache=shared",
			want: "file:app?mode=memory&cache=shared",
		},
		{
			// The "?" belongs to a directory name, so the path has no query
			// and the options open one.
			name: "question mark in a directory name",
			dsn:  "/tmp/a?b/app.db",
			want: "/tmp/a?b/app.db?_busy_timeout=5000&_journal_mode=wal",
		},
		{
			name: "question mark in a directory name and a real query",
			dsn:  "/tmp/a?b/app.db?cache=shared",
			want: "/tmp/a?b/app.db?cache=shared&_busy_timeout=5000&_journal_mode=wal",
		},
		{
			name: "absolute path",
			dsn:  "/var/lib/app/app.db",
			want: "/var/lib/app/app.db?_busy_timeout=5000&_journal_mode=wal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cleanDSN("sqlite", tt.dsn))
		})
	}
}
