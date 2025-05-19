package sfrouting

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.snappfood.ir/backend/go/packages/sf-routing/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.elastic.co/apm/module/apmgin/v2"
)

// =============================================================================
// Core Interfaces
// =============================================================================

// Healthy defines the interface for health check services
type Healthy interface {
	Health(ctx context.Context) error
}

// PrometheusExporter defines the interface for Prometheus exporter
type PrometheusExporter interface {
	// Handler returns an http.Handler that will handle Prometheus metrics requests
	Handler() http.Handler
}

// Router defines the interface for registering routes
type Router interface {
	Routes(router *gin.Engine)
}

// RouterGroup defines the interface for registering routes in a group
type RouterGroup interface {
	Routes(router *gin.RouterGroup)
}

// Middleware is a function that takes a gin.HandlerFunc and returns a gin.HandlerFunc
type Middleware = middlewares.Middleware

// ErrorHandler is a function that handles errors
type ErrorHandler func(c *gin.Context, err error)

// =============================================================================
// Configuration Types
// =============================================================================

// CorsConfig holds CORS configuration settings
type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// PrometheusConfig holds Prometheus configuration settings
type PrometheusConfig struct {
	Enabled bool
	Path    string
}

// DefaultPrometheusConfig returns the default Prometheus configuration
func DefaultPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		Enabled: false,
		Path:    "/metrics",
	}
}

// SwaggerConfig holds all Swagger configuration settings
type SwaggerConfig struct {
	Enabled  bool
	Title    string
	Version  string
	Host     string
	BasePath string
	Schemes  []string
}

// DefaultSwaggerConfig returns the default Swagger configuration
func DefaultSwaggerConfig() SwaggerConfig {
	return SwaggerConfig{
		Enabled:  false,
		Title:    "API Documentation",
		Version:  "1.0",
		Host:     "",
		BasePath: "/",
		Schemes:  []string{"http", "https"},
	}
}

// =============================================================================
// Registry Types and Global Instance
// =============================================================================

// Option is a function to customize the Registry
type Option func(*Registry)

// Registry is a global registry for Gin configuration
type Registry struct {
	mu                 sync.RWMutex
	engine             *gin.Engine
	healthChecks       []Healthy
	globalMiddlewares  []Middleware
	errorHandler       ErrorHandler
	logger             Logger
	SwaggerEnabled     bool
	SwaggerTitle       string
	SwaggerVersion     string
	SwaggerHost        string
	SwaggerBasePath    string
	SwaggerSchemes     []string
	corsConfig         CorsConfig
	corsApplied        bool
	prometheusExporter PrometheusExporter
	prometheusConfig   PrometheusConfig
}

// Global registry instance
var globalRegistry = &Registry{
	engine:            gin.Default(),
	healthChecks:      make([]Healthy, 0),
	globalMiddlewares: make([]Middleware, 0),
	logger:            nil,
	errorHandler:      nil,
	SwaggerEnabled:    false,
	SwaggerTitle:      "API Documentation",
	SwaggerVersion:    "1.0",
	SwaggerHost:       "",
	SwaggerBasePath:   "/",
	SwaggerSchemes:    []string{"http", "https"},
	corsConfig:        CorsConfig{},
	prometheusConfig:  PrometheusConfig{Enabled: false, Path: "/metrics"},
}

// =============================================================================
// Registry Configuration Options
// =============================================================================

// WithHealthChecks adds health check services to the registry
func WithHealthChecks(checkers ...Healthy) Option {
	return func(r *Registry) {
		r.healthChecks = append(r.healthChecks, checkers...)
	}
}

// WithPrometheusExporter sets the Prometheus exporter and configuration
func WithPrometheusExporter(exporter PrometheusExporter, config PrometheusConfig) Option {
	return func(r *Registry) {
		r.prometheusExporter = exporter
		r.prometheusConfig = config
	}
}

// WithLogger sets a custom logger for the registry
func WithLogger(logger Logger) Option {
	return func(r *Registry) {
		r.logger = logger
	}
}

// WithSwagger configures all Swagger settings in one place
func WithSwagger(config SwaggerConfig) Option {
	return func(r *Registry) {
		r.SwaggerEnabled = config.Enabled
		r.SwaggerTitle = config.Title
		r.SwaggerVersion = config.Version
		r.SwaggerHost = config.Host
		r.SwaggerBasePath = config.BasePath
		r.SwaggerSchemes = config.Schemes
	}
}

// WithErrorHandler sets a custom error handler
func WithErrorHandler(handler ErrorHandler) Option {
	return func(r *Registry) {
		r.errorHandler = handler
	}
}

// WithGlobalMiddleware adds a global middleware
func WithGlobalMiddleware(middleware Middleware) Option {
	return func(r *Registry) {
		r.globalMiddlewares = append(r.globalMiddlewares, middleware)
	}
}

// WithGinConfig allows direct configuration of the Gin engine
func WithGinConfig(config func(*gin.Engine)) Option {
	return func(r *Registry) {
		config(r.engine)
	}
}

// WithCorsConfig configures CORS settings for the Gin engine
func WithCorsConfig(config CorsConfig) Option {
	return func(r *Registry) {
		// Store the CORS configuration
		r.corsConfig = config

		// Apply CORS middleware to the Gin engine
		r.engine.Use(func(c *gin.Context) {
			// Set CORS headers
			if len(config.AllowOrigins) > 0 {
				c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(config.AllowOrigins, ", "))
			}

			if len(config.AllowMethods) > 0 {
				c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
			}

			if len(config.AllowHeaders) > 0 {
				c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
			}

			if len(config.ExposeHeaders) > 0 {
				c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
			}

			if config.AllowCredentials {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if config.MaxAge > 0 {
				c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", int(config.MaxAge.Seconds())))
			}

			// Handle preflight requests
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		})
	}
}

// =============================================================================
// Connection Management
// =============================================================================

// RegisterConnection configures the Gin engine with the provided options
func RegisterConnection(opts ...Option) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// Apply APM middleware
	globalRegistry.engine.Use(apmgin.Middleware(globalRegistry.engine))

	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	// Log startup if logger is provided
	if globalRegistry.logger != nil {
		extraMap := map[string]interface{}{}

		globalRegistry.logger.InfoWithCategory(
			Category.System.General,
			SubCategory.Operation.Startup,
			"Service registry initialized",
			extraMap,
		)
	}

	// Apply global middlewares
	for _, middleware := range globalRegistry.globalMiddlewares {
		globalRegistry.engine.Use(func(c *gin.Context) {
			// Create a handler function that calls c.Next()
			nextHandler := func(c *gin.Context) {
				c.Next()
			}
			middleware(nextHandler)(c)
		})
	}

	// Register health check endpoint
	globalRegistry.engine.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		var errs []string

		for _, check := range globalRegistry.healthChecks {
			if err := check.Health(ctx); err != nil {
				errs = append(errs, err.Error())
			}
		}

		if len(errs) > 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"errors": errs,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Register Prometheus exporter if enabled
	if globalRegistry.prometheusExporter != nil && globalRegistry.prometheusConfig.Enabled {
		handler := globalRegistry.prometheusExporter.Handler()

		// Register the Prometheus handler with Gin
		globalRegistry.engine.GET(globalRegistry.prometheusConfig.Path, func(c *gin.Context) {
			handler.ServeHTTP(c.Writer, c.Request)
		})

		if globalRegistry.logger != nil {
			extraMap := map[string]interface{}{
				ExtraKey.HTTP.Path: globalRegistry.prometheusConfig.Path,
			}

			globalRegistry.logger.InfoWithCategory(
				Category.Monitoring.Metrics,
				SubCategory.Operation.Configuration,
				"Prometheus metrics endpoint registered",
				extraMap,
			)
		}
	}

	// Register Swagger if enabled
	if globalRegistry.SwaggerEnabled {
		// Default to Swagger UI
		url := ginSwagger.URL("/swagger/doc.json") // The URL pointing to API definition
		globalRegistry.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

		if globalRegistry.logger != nil {
			extraMap := map[string]interface{}{
				ExtraKey.HTTP.Path: "/swagger/*any",
			}

			globalRegistry.logger.InfoWithCategory(
				Category.API.Documentation,
				SubCategory.Operation.Configuration,
				"Swagger API documentation endpoint registered",
				extraMap,
			)
		}
	}

	// Set custom error handler if provided
	if globalRegistry.errorHandler != nil {
		globalRegistry.engine.Use(func(c *gin.Context) {
			c.Next()
			if len(c.Errors) > 0 {
				globalRegistry.errorHandler(c, c.Errors.Last().Err)
			}
		})
	}

	return nil
}

// =============================================================================
// Public API
// =============================================================================

// RegisterHealthCheck registers a health check service
// This method is kept for backward compatibility
func RegisterHealthCheck(checker ...Healthy) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	globalRegistry.healthChecks = append(globalRegistry.healthChecks, checker...)
}

// RegisterRouter registers a router
func RegisterRouter(router Router) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	router.Routes(globalRegistry.engine)
}

// RegisterRouterGroup registers a router group
func RegisterRouterGroup(groupPath string, router RouterGroup) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	group := globalRegistry.engine.Group(groupPath)
	router.Routes(group)
}

// StartServer starts the HTTP server
func StartServer(addr string) error {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	if globalRegistry.logger != nil {
		extraMap := map[string]interface{}{
			ExtraKey.Network.HostIP: addr,
		}

		globalRegistry.logger.InfoWithCategory(
			Category.Infrastructure.Network,
			SubCategory.Operation.Startup,
			"Starting server",
			extraMap,
		)
	}
	return globalRegistry.engine.Run(addr)
}

// GetEngine returns the Gin engine
func GetEngine() *gin.Engine {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	return globalRegistry.engine
}

// DefaultErrorHandler is a default error handler
func DefaultErrorHandler(c *gin.Context, err error) {
	if globalRegistry.logger != nil {
		extraMap := map[string]interface{}{
			ExtraKey.Error.ErrorMessage: err.Error(),
		}

		globalRegistry.logger.ErrorWithCategory(
			Category.Error.Error,
			SubCategory.Status.Error,
			"Request error",
			extraMap,
		)
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
