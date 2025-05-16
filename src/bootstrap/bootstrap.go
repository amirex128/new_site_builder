package bootstrap

import (
	"context"
	"go-boilerplate/src/bootstrap/service_provider"
)

func GrpcBootstrap(ctx context.Context) *Container {

	logger := serviceprovider.LoggerProvider(ctx)

	cfg := serviceprovider.ConfigProvider(logger)

	serviceprovider.ExternalServicesProvider(cfg, logger)

	serviceprovider.ServiceDiscoveryProvider(cfg, logger)

	serviceprovider.GrpcRouterProvider()

	container := ContainerProvider(ctx, cfg, logger)

	serviceprovider.HttpProvider(logger)

	serviceprovider.GrpcProvider(logger)

	serviceprovider.MemoryLoaderProvider(logger)

	return container
}
func ConsumerBootstrap(ctx context.Context) *Container {

	logger := serviceprovider.LoggerProvider(ctx)

	cfg := serviceprovider.ConfigProvider(logger)

	serviceprovider.ExternalServicesProvider(cfg, logger)

	serviceprovider.ServiceDiscoveryProvider(cfg, logger)

	serviceprovider.RabbitProvider(logger)

	serviceprovider.KafkaProvider()

	container := ContainerProvider(ctx, cfg, logger)

	serviceprovider.HttpProvider(logger)

	serviceprovider.GrpcProvider(logger)

	serviceprovider.MemoryLoaderProvider(logger)

	return container
}
func HttpBootstrap(ctx context.Context) *Container {

	logger := serviceprovider.LoggerProvider(ctx)

	cfg := serviceprovider.ConfigProvider(logger)

	serviceprovider.ExternalServicesProvider(cfg, logger)

	serviceprovider.ServiceDiscoveryProvider(cfg, logger)

	serviceprovider.RouterProvider(logger)

	container := ContainerProvider(ctx, cfg, logger)

	serviceprovider.HttpProvider(logger)

	serviceprovider.GrpcProvider(logger)

	serviceprovider.MemoryLoaderProvider(logger)

	return container
}
