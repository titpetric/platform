package platform

import (
	"context"
)

func Start() error {
	svc, err := StartPlatform(context.Background(), nil)
	if err != nil {
		return err
	}

	svc.Wait()
	return nil
}

func StartPlatform(ctx context.Context, options *Options) (*Platform, error) {
	svc, err := NewPlatform(options)
	if err != nil {
		return nil, err
	}

	svc.Serve(ctx)
	return svc, nil
}
