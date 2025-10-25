package platform

import (
	"context"
)

func Start() error {
	svc, err := StartPlatform(context.Background())
	if err != nil {
		return err
	}

	svc.Wait()
	return nil
}

func StartPlatform(ctx context.Context, opts ...*Options) (*Platform, error) {
	svc, err := NewPlatform(opts...)
	if err != nil {
		return nil, err
	}

	svc.Serve(ctx)
	return svc, nil
}
