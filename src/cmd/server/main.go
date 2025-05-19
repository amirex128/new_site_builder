package main

import (
	"context"
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/amirex128/new_site_builder/src/internal/api/router"
)

func main() {
	ctx := context.Background()

	container := bootstrap.HttpBootstrap(ctx)

	handlers := bootstrap.HandlerBootstrap(container)

	router.InitServer(handlers, container, container.Logger, container.Config)
}
