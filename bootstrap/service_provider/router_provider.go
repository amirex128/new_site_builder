package serviceprovider

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	"github.com/amirex128/new_site_builder/bootstrap/exporter"
	"github.com/amirex128/new_site_builder/bootstrap/healthcheck"
	"github.com/gin-gonic/gin"
)

func RouterProvider(logger sflogger.Logger) {
	err := sfrouting.RegisterConnection(
		sfrouting.WithLogger(logger),
		sfrouting.WithHealthChecks(&healthcheck.BaseHealthCheck{}),
		sfrouting.WithPrometheusExporter(&exporter.BasePrometheusExporter{}, sfrouting.PrometheusConfig{
			Enabled: true,
			Path:    "/metrics",
		}),
		sfrouting.WithSwagger(sfrouting.SwaggerConfig{
			Enabled:  true,
			Title:    "SF-Routing API",
			Version:  "1.0",
			BasePath: "/api/v1",
			Schemes:  []string{"http", "https"},
		}),
		sfrouting.WithCorsConfig(sfrouting.CorsConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			AllowCredentials: true,
		}),
		sfrouting.WithGinConfig(func(engine *gin.Engine) {
			gin.SetMode(gin.ReleaseMode)
		}),
	)

	if err != nil {
		logger.ErrorWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, fmt.Sprintf("Failed to register routers: %v", err), nil)
	}
	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded Router", nil)

}
