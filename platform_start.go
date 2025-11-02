package platform

import (
	"context"
)

func Start(ctx context.Context, options *Options) (*Platform, error) {
	svc, err := New(options)
	if err != nil {
		return nil, err
	}

	if err := svc.Start(ctx); err != nil {
		return nil, err
	}
	return svc, nil
}
