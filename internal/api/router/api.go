package router

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	"github.com/amirex128/new_site_builder/bootstrap"
	"github.com/amirex128/new_site_builder/docs"
	middleware2 "github.com/amirex128/new_site_builder/internal/api/middleware"
	contract2 "github.com/amirex128/new_site_builder/internal/contract"
)

func RunServer(handlers *bootstrap.HandlerManager, container contract2.IContainer, logger sflogger.Logger, cnf contract2.IConfig) {
	sfrouting.AddGlobalMiddlewares(middleware2.LoggerMiddleware(logger), middleware2.ErrorHandlerMiddleware(logger))
	RegisterRoutes(handlers, container)
	docs.SwaggerInfo.BasePath = "/api/v1"
	err := sfrouting.StartServer(fmt.Sprintf(":%s", cnf.GetString("app_port")))
	if err != nil {
		logger.Errorf("Failed to start server: %s", err.Error())
		return
	}
}

// RegisterRoutes registers all routes
func RegisterRoutes(h *bootstrap.HandlerManager, container contract2.IContainer) {
	// Register user routes in a group
	sfrouting.RegisterRouterGroup("/api/v1", &RouterV1{h: h, container: container})
}
