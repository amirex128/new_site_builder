package main

import (
	"context"
	bootstrap2 "github.com/amirex128/new_site_builder/bootstrap"
	"github.com/amirex128/new_site_builder/internal/api/router"
)

// @title My API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8585
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @security BearerAuth
func main() {
	ctx := context.Background()

	container := bootstrap2.HttpServerBootstrap(ctx)

	handlers := bootstrap2.HttpHandlerBootstrap(container)

	router.RunServer(handlers, container, container.Logger, container.Config)
}
