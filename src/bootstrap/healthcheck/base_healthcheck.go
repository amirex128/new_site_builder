package healthcheck

import "context"

type BaseHealthCheck struct {
}

func (BaseHealthCheck) Health(ctx context.Context) error {

	return nil
}
