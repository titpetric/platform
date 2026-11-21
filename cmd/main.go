package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/titpetric/oida"

	"github.com/titpetric/platform"
)

// Main runs the platform server with a manager and blocks until shutdown.
// A caught SIGHUP reloads, while SIGINT or SIGTERM stops it. A variadic parameter injects options from a test.
func Main(ctx context.Context, options ...*platform.Options) {
	var option *platform.Options
	for _, opt := range options {
		option = opt
		break
	}

	m := platform.NewManager(option)
	if err := m.Start(ctx); err != nil {
		err = fmt.Errorf("exit error: %w", err)
		oida.RecordError(ctx, err)

		// The platform never came up, so there is no platform logger to
		// write this through.
		slog.Error("platform start failed", "error", err)
		os.Exit(1)
	}

	m.Wait()
}
