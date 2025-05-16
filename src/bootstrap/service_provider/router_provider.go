package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	"git.snappfood.ir/backend/go/packages/sf-routing/middlewares"
	"github.com/gin-gonic/gin"
	"go-boilerplate/src/bootstrap/exporter"
	"go-boilerplate/src/bootstrap/healthcheck"
	"log"
)

func RouterProvider(logger sflogger.Logger) {
	// Configure the server
	err := sfrouting.RegisterConnection(
		sfrouting.WithLogger(logger),
		sfrouting.WithHealthChecks(&healthcheck.BaseHealthCheck{}),
		sfrouting.WithPrometheusExporter(&exporter.BasePrometheusExporter{}, sfrouting.PrometheusConfig{
			Enabled: true,
			Path:    "/metrics",
		}),
		sfrouting.WithSwagger(sfrouting.SwaggerConfig{
			Enabled:  true,
			Title:    "SF-Routing Example API",
			Version:  "1.0",
			Host:     "localhost:8080",
			BasePath: "/",
			Path:     "/swagger/*any",
			Schemes:  []string{"http", "https"},
		}),
		sfrouting.WithErrorHandler(sfrouting.DefaultErrorHandler),
		sfrouting.WithGlobalMiddleware(middlewares.LoggedinMiddle),
		sfrouting.WithCorsConfig(sfrouting.CorsConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			AllowCredentials: true,
		}),
		sfrouting.WithGinConfig(func(engine *gin.Engine) {
			// Set Gin to release mode
			gin.SetMode(gin.ReleaseMode)

			// Configure custom recovery middleware
			engine.Use(gin.Recovery())

			// Configure custom static files
			engine.Static("/static", "./static")

			// Add a custom middleware to all routes
			engine.Use(func(c *gin.Context) {
				c.Set("custom_middleware", "applied")
				c.Next()
			})
		}),
	)

	if err != nil {
		log.Fatalf("Failed to register Routers: %v", err)
	}

}
