package platform

import (
	"context"
)

func Start(ctx context.Context, options *Options) (*Platform, error) {
	svc := New(options)
	if err := svc.Start(ctx); err != nil {
		return nil, err
	}
	return svc, nil
}
