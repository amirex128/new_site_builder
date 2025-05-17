package contract

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	"github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
)

type IContainer interface {
	GetConfig() IConfig
	GetSimpleRepo() repository.ISampleRepository
	GetSimpleZoodFoodRepo() repository.ISampleRepository
	GetFoodPartyCash() cache.ICacheService
	GetStockCacheTransient() cache.ICacheService
	GetLogger() sflogger.Logger
}
