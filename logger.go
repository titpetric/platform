package platform

import (
	"context"
	"log/slog"
)

// Logger is the interface the platform writes output through: the subset of
// *slog.Logger it needs. A *slog.Logger assigns to Platform.Logger as it is.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// discard swallows platform output, backing Options.Quiet.
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

// loggerFromContext returns the platform logger bound to the context.
func loggerFromContext(ctx context.Context) Logger {
	if p := FromContext(ctx); p != nil {
		return p.logger()
	}
	return discard
}
