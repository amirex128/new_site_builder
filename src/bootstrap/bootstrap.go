package bootstrap

import (
	"context"
	"github.com/amirex128/new_site_builder/src/bootstrap/service_provider"
)

func HttpBootstrap(ctx context.Context) *Container {

	logger := serviceprovider.LoggerProvider(ctx)

	cfg := serviceprovider.ConfigProvider(logger)

	serviceprovider.MongoProvider(cfg, logger)

	serviceprovider.RedisProvider(cfg, logger)

	serviceprovider.MysqlProvider(cfg, logger)

	serviceprovider.ServiceDiscoveryProvider(cfg, logger)

	serviceprovider.RouterProvider(logger)

	container := ContainerProvider(ctx, cfg, logger)

	serviceprovider.HttpProvider(logger)

	serviceprovider.GrpcProvider(logger)

	serviceprovider.MemoryLoaderProvider(logger)

	return container
}
