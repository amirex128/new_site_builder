package bootstrap

import (
	"context"
	"git.snappfood.ir/backend/go/packages/sf-logger"
	"git.snappfood.ir/backend/go/packages/sf-orm"
	"git.snappfood.ir/backend/go/packages/sf-redis"
	"github.com/amirex128/new_site_builder/src/config"
	"github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
	"github.com/amirex128/new_site_builder/src/internal/infra/repository/mysql"
	"github.com/amirex128/new_site_builder/src/internal/infra/service"
)

func ContainerProvider(ctx context.Context, cfg *config.Config, logger sflogger.Logger) *Container {

	return &Container{
		Config: cfg,
		Logger: logger,

		//todo: create name constant
		FoodPartyCache: service.NewRedis(sfredis.MustClient(ctx, "foodparty")),

		// for transient
		stockCacheTransient: func() cache.ICacheService {
			return service.NewRedis(sfredis.MustClient(ctx, "stock"))
		},

		SampleRepo:         mysql.NewSampleRepository(sform.MustDB("search_db")),
		SampleZoodFoodRepo: mysql.NewSampleRepository(sform.MustDB("zoodfood_db")),
	}
}
