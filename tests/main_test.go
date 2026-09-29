// Package platform_test holds the black box tests of the platform module.
//
// It is a separate Go module so that the sql drivers it registers stay out
// of the platform dependency tree. A driver is the choice of whoever builds
// the binary: platform speaks to database/sql and never imports one. Pinning
// a driver here would push that pin onto every consumer and duplicate the
// dependency tree of anyone already on a newer version.
package platform_test

import (
	"testing"

	// The drivers under test. Blank imports register them with
	// database/sql, which is the only thing platform needs of them.
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// TestMain runs the package tests. The driver registrations above happen in
// package init, before it is called, so every test in this package can open
// a mysql, pgx or sqlite DSN without importing a driver itself.
func TestMain(m *testing.M) {
	m.Run()
}
