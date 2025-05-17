package serviceprovider

import (
	"github.com/amirex128/new_site_builder/src/config"

	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func ExternalServicesProvider(cfg *config.Config, logger sflogger.Logger) {
	//mysqlConnector := orm.GetConnector()
	//redisConnector := redis.GetConnector()
	//elasticConnector := elastic.GetConnector()
	//
	//logger.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Operation.Registration, "Registering external services", nil)
	//
	//err := serviceregistry.RegisterServices(
	//	"server",
	//	"src/config/config.json",
	//	cfg,
	//	mysqlConnector,
	//	redisConnector,
	//	elasticConnector,
	//	nil,
	//	nil,
	//	serviceregistry.WithLogger(logger),
	//)
	//
	//if err != nil {
	//	extraMap := map[string]interface{}{
	//		sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
	//	}
	//	logger.FatalWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Error, "Failed to register service", extraMap)
	//}

	ElasticProvider(logger)
	RedisProvider(logger)
	RabbitProvider(logger)
	MysqlProvider(logger)

	logger.InfoWithCategory(sflogger.Category.Service.External, sflogger.SubCategory.Status.Success, "External services registered successfully", nil)
}
