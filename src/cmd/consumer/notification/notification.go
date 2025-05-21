package main

import (
	"context"
	"github.com/amirex128/new_site_builder/src/bootstrap"
	consumerrouter "github.com/amirex128/new_site_builder/src/internal/api/consumer_router"
)

func main() {
	ctx := context.Background()

	container := bootstrap.ConsumerServerBootstrap(ctx)

	handlers := bootstrap.ConsumerHandlerBootstrap(container)

	consumerrouter.RunServer(handlers, container, container.Logger, container.Config)
}
