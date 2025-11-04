package cmd

import (
	"context"
	"fmt"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/pkg/telemetry"

	// Add platform drivers and modules.
	_ "github.com/titpetric/platform/pkg/drivers"
)

// Main is the entrypoint for the app.
//
// It's expected to have control of the app lifecycle. An application
// exit is not expected to be graceful in case of errors. Main starts
// the platform server with modules loaded beforehand. It is blocking
// until server shutdown from cancellation of the context, or a caught
// SIGTERM, an OS control signal hinting the app should exit.
//
// The variadic parameter allows to inject options from test.
func Main(ctx context.Context, options ...*platform.Options) {
	var option *platform.Options
	for _, opt := range options {
		option = opt
		break
	}

	p, err := platform.Start(ctx, option)
	if err != nil {
		telemetry.Fatal(ctx, fmt.Errorf("exit error: %w", err))
	}

	p.Wait()
}
