package router

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/docs"
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/amirex128/new_site_builder/src/internal/api/middleware"
	"github.com/amirex128/new_site_builder/src/internal/contract"

	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
)

func RunServer(handlers *bootstrap.HandlerManager, container contract.IContainer, logger sflogger.Logger, cnf contract.IConfig) {
	sfrouting.AddGlobalMiddlewares(middleware.LoggerMiddleware(logger), middleware.ErrorHandlerMiddleware(logger))
	RegisterRoutes(handlers, container)
	docs.SwaggerInfo.BasePath = "/api/v1"
	err := sfrouting.StartServer(fmt.Sprintf(":%s", cnf.GetString("app_port")))
	if err != nil {
		logger.Errorf("Failed to start server: %s", err.Error())
		return
	}
}

// RegisterRoutes registers all routes
func RegisterRoutes(h *bootstrap.HandlerManager, container contract.IContainer) {
	// Register user routes in a group
	sfrouting.RegisterRouterGroup("/api/v1", &RouterV1{h: h, container: container})
}
