package platform

import (
	"context"
	"log/slog"
	"testing"

	"github.com/titpetric/platform/internal/assert"
)

// TestLogger covers the logger a platform and a manager return, which is never nil.
func TestLogger(t *testing.T) {
	var logger Logger = slog.Default()

	p := New(NewTestOptions())
	p.Logger = logger
	assert.Equal(t, logger, p.logger())

	p.Logger = nil
	assert.Equal(t, discard, p.logger())

	m := NewManager(NewTestOptions())
	m.Logger = logger
	assert.Equal(t, logger, m.logger())

	m.Logger = nil
	assert.Equal(t, discard, m.logger())
}

// TestLoggerFromContext covers the lookup a registry falls back on when no
// platform is bound to the context.
func TestLoggerFromContext(t *testing.T) {
	assert.Equal(t, discard, loggerFromContext(context.Background()))

	p := New(NewTestOptions())
	ctx := platformContext.SetContext(context.Background(), p)
	assert.Equal(t, p.Logger, loggerFromContext(ctx))
}
