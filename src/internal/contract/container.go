package contract

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"go-boilerplate/src/internal/contract/repository"
	"go-boilerplate/src/internal/contract/service/cache"
)

type IContainer interface {
	GetConfig() IConfig
	GetSimpleRepo() repository.ISampleRepository
	GetSimpleZoodFoodRepo() repository.ISampleRepository
	GetFoodPartyCash() cache.ICacheService
	GetStockCacheTransient() cache.ICacheService
	GetLogger() sflogger.Logger
}
