package main

import (
	"context"
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/amirex128/new_site_builder/src/internal/api/router"
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

	container := bootstrap.HttpBootstrap(ctx)

	handlers := bootstrap.HandlerBootstrap(container)

	router.InitServer(handlers, container, container.Logger, container.Config)
}
