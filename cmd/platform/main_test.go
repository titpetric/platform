package main

import (
	"os"
	"testing"
)

// TestMain runs the package tests. main blocks on the platform manager, so
// nothing here calls it without a server to serve.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
