package platform

import "log/slog"

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
