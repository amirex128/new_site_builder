package bootstrap

import (
	"context"
	"github.com/amirex128/new_site_builder/src/bootstrap/service_provider"
)

func HttpServerBootstrap(ctx *context.Context) *Container {

	logger := serviceprovider.LoggerProvider()

	cfg := serviceprovider.ConfigProvider(logger)

	serviceprovider.MongoProvider(cfg, logger)

	serviceprovider.RedisProvider(cfg, logger)

	serviceprovider.MysqlProvider(cfg, logger)

	serviceprovider.RabbitProvider(cfg, logger)

	serviceprovider.ServiceDiscoveryProvider(cfg, logger)

	serviceprovider.RouterProvider(logger)

	container := ContainerProvider(ctx, cfg, logger)

	serviceprovider.HttpRequestProvider(logger)

	serviceprovider.GrpcRequestProvider(logger)

	serviceprovider.MemoryLoaderProvider(logger)

	return container
}

func ConsumerServerBootstrap(ctx *context.Context) *Container {

	logger := serviceprovider.LoggerProvider()

	cfg := serviceprovider.ConfigProvider(logger)

	serviceprovider.MongoProvider(cfg, logger)

	serviceprovider.RedisProvider(cfg, logger)

	serviceprovider.MysqlProvider(cfg, logger)

	serviceprovider.RabbitProvider(cfg, logger)

	serviceprovider.ServiceDiscoveryProvider(cfg, logger)

	serviceprovider.RouterProvider(logger)

	container := ContainerProvider(ctx, cfg, logger)

	serviceprovider.HttpRequestProvider(logger)

	serviceprovider.GrpcRequestProvider(logger)

	serviceprovider.MemoryLoaderProvider(logger)

	return container
}
