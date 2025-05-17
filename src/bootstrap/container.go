package bootstrap

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/internal/contract"
	"github.com/amirex128/new_site_builder/src/internal/contract/repository"
	cache2 "github.com/amirex128/new_site_builder/src/internal/contract/service/cache"
)

// Container
type Container struct {
	Config             contract.IConfig
	SampleRepo         repository.ISampleRepository
	SampleZoodFoodRepo repository.ISampleRepository
	MemoryLoader       cache2.IMemoryLoader
	FoodPartyCache     cache2.ICacheService

	stockCacheTransient func() cache2.ICacheService

	Logger sflogger.Logger
}

func (c *Container) GetConfig() contract.IConfig {
	return c.Config
}
func (c *Container) GetSimpleRepo() repository.ISampleRepository {
	return c.SampleRepo
}

func (c *Container) GetSimpleZoodFoodRepo() repository.ISampleRepository {
	return c.SampleZoodFoodRepo
}

func (c *Container) GetFoodPartyCash() cache2.ICacheService {
	return c.FoodPartyCache
}

func (c *Container) GetStockCacheTransient() cache2.ICacheService {
	return c.stockCacheTransient()
}

func (c *Container) GetLogger() sflogger.Logger {
	return c.Logger
}
