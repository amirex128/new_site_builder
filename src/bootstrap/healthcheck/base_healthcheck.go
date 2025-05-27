package healthcheck

import (
	"context"
	sfelastic "git.snappfood.ir/backend/go/packages/sf-elasticsearch-client/v8"
	sfmongo "git.snappfood.ir/backend/go/packages/sf-mongo"
	sform "git.snappfood.ir/backend/go/packages/sf-orm"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

type BaseHealthCheck struct {
}

func (BaseHealthCheck) Health(ctx context.Context) error {

	err := sfredis.Health(ctx)
	if err != nil {
		return err
	}
	err = sfelastic.Health(ctx)
	if err != nil {
		return err
	}
	err = sform.Health(ctx)
	if err != nil {
		return err
	}
	err = sfmongo.Health(ctx)
	if err != nil {
		return err
	}
	err = sfrabbitmq.Health(ctx)
	if err != nil {
		return err
	}

	return nil
}
