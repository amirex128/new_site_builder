package main

import (
	"context"
	"go-boilerplate/src/bootstrap"
	"go-boilerplate/src/internal/api/router"
)

func main() {
	ctx := context.Background()

	container := bootstrap.HttpBootstrap(ctx)

	handlers := bootstrap.HandlerBootstrap(container)

	router.InitServer(handlers)
}
