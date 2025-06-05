package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/config"
)

func ServiceDiscoveryProvider(cfg *config.Config, logger sflogger.Logger) {
	// Initialize and register service with the service registry
	//serviceConfig := serviceregistry.ServiceRegistryConfig{
	//	// Registry configuration
	//	Registry: serviceregistry.RegistryConfig{
	//		Type:     "consul",
	//		Address:  cfg.ServerServiceRegistryAddress,
	//		Token:    cfg.ServerServiceRegistryToken,
	//		Username: cfg.ServerServiceRegistryUsername,
	//		Password: cfg.ServerServiceRegistryPassword,
	//	},
	//
	//	// Service configuration
	//	Service: serviceregistry.ServiceConfig{
	//		Name:    "search",
	//		Port:    cfg.ServerPort,
	//		Address: cfg.ServerServiceAddress,
	//		Tags:    "search",
	//		Scheme:  cfg.ServerServiceScheme,
	//		HealthCheck: serviceregistry.HealthCheck{
	//			URL:      cfg.ServerServiceHealthCheckURL,
	//			TTL:      cfg.ServerServiceHealthCheckTTL,
	//			Interval: cfg.ServerServiceHealthCheckInterval,
	//		},
	//	},
	//
	//	// Behavior configuration
	//	SetupGracefulShutdown: true,
	//}
	//
	//serviceID, err := serviceregistry.InitAndRegisterService(logger, serviceConfig)
	//if err != nil {
	//	extraMap := map[string]interface{}{
	//		sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
	//		"address":                            cfg.ServerServiceRegistryAddress,
	//	}
	//	logger.ErrorWithCategory(sflogger.Category.System.Internal, sflogger.SubCategory.Status.Error, "Failed to initialize service registry", extraMap)
	//}
	//
	//// Log the service ID if registration was successful
	//if serviceID != "" {
	//	extraMap := map[string]interface{}{
	//		"serviceID": serviceID,
	//		"type":      "consul",
	//	}
	//	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Service registered successfully", extraMap)
	//}
}
