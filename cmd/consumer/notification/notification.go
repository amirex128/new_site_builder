package main

import (
	"context"
	bootstrap2 "github.com/amirex128/new_site_builder/bootstrap"
	"github.com/amirex128/new_site_builder/internal/api/consumer_router"
)

func main() {
	ctx := context.Background()

	container := bootstrap2.ConsumerServerBootstrap(ctx)

	handlers := bootstrap2.ConsumerHandlerBootstrap(container)

	consumerrouter.RunServer(&ctx, handlers, container)
}
