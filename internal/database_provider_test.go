package internal

import (
	"database/sql"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestNewDatabaseProvider covers a provider with no credentials.
func TestNewDatabaseProvider(t *testing.T) {
	var provider *DatabaseProvider = NewDatabaseProvider(func(string, string) (*sql.DB, error) {
		return nil, nil
	})

	assert.NotNil(t, provider)
	assert.Empty(t, provider.List())

	provider.Register("main", "sqlite://:memory:")
	assert.Equal(t, []string{"main"}, provider.List())
}
