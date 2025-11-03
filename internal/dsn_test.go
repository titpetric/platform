package internal

import (
	"testing"

	"github.com/titpetric/platform/pkg/require"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanDSN(tt.dsn)
			require.Equal(t, tt.want, got)
		})
	}
}
