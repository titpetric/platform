package platform

import (
	"context"

	"github.com/titpetric/platform/internal/tracing"
)

func CaptureError(ctx context.Context, err error) {
	tracing.CaptureError(ctx, err)
}
