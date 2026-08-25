package platform

import (
	"context"
	"log/slog"
)

// Logger is the interface the platform writes its own output through.
// It is the subset of *slog.Logger the platform needs, so a *slog.Logger
// can be assigned to Platform.Logger as it is.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// discard swallows platform output. It backs Options.Quiet, and stands in
// when Platform.Logger has been set to nil.
var discard Logger = slog.New(slog.DiscardHandler)

// logger returns the platform logger, never nil.
func (p *Platform) logger() Logger {
	if p.Logger == nil {
		return discard
	}
	return p.Logger
}

// logger returns the manager logger, never nil.
func (m *Manager) logger() Logger {
	if m.Logger == nil {
		return discard
	}
	return m.Logger
}

// loggerFromContext returns the logger of the platform bound to the context.
// A Registry is usable as a bare value, outside of a platform, so there may
// be no platform in the context, in which case output is discarded.
func loggerFromContext(ctx context.Context) Logger {
	if p := FromContext(ctx); p != nil {
		return p.logger()
	}
	return discard
}
