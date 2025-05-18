package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/config"
)

func ExternalServicesProvider(cfg *config.Config, logger sflogger.Logger) {
	//ElasticProvider(logger)
	RedisProvider(cfg, logger)
	//RabbitProvider(logger)
	MysqlProvider(cfg, logger)

	logger.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "External services registered successfully", nil)
}
