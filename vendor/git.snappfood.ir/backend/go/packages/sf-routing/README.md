# SF-Routing

SF-Routing is a Golang library that wraps the Gin framework to make it easy to use in projects with just a few lines of code. It provides features for health checks, middleware management, error handling, Swagger documentation, and more.

## Features

- **Health Check Registration**: Register health check services that implement the `Healthy` interface.
- **Route Registration**: Register routes with middleware support.
- **Router Groups**: Organize routes into groups for better structure.
- **Middleware Management**: Add global and per-route middleware.
- **Error Handling**: Set a custom error handler for all routes.
- **API Documentation**: Generate OpenAPI documentation with Swagger UI.
- **Custom Logger**: Use a custom logger that implements the `Logger` interface.
- **Direct Gin Configuration**: Configure the Gin engine directly with custom settings.
- **CORS Configuration**: Easily set up Cross-Origin Resource Sharing (CORS) for your API.
- **Prometheus Integration**: Expose Prometheus metrics by implementing the `PrometheusExporter` interface.

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-routing
```

# OpenAPI Documentation

SF-Routing supports OpenAPI documentation using Swagger UI.

## Prerequisites

To use the OpenAPI documentation features, you need to install the following packages:

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## Using Swagger UI

To use Swagger UI in your application:

```go
err := sfrouting.RegisterConnection(
	sfrouting.WithSwagger(sfrouting.SwaggerConfig{
		Enabled:  true,
		Title:    "SF-Routing API",
		Version:  "1.0",
		Host:     "localhost:8080",
		BasePath: "/api/v1",
		Schemes:  []string{"http", "https"},
	}),
)
```

## Documenting Your API

Add Swagger annotations to your handler functions:

```go
// CreateUser godoc
// @Summary      Create a user
// @Description  create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User object"
// @Success      201   {object}  User
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users [post]
// @Security     BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Implementation
}
```

## Generating Documentation

Generate the documentation using the `swag` command:

```bash
# Exit on error
set -e

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "swag is not installed. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Create docs directory if it doesn't exist
mkdir -p docs

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g path/to/api.go -o ./docs

echo "Documentation generated successfully!"
echo "You can access the Swagger UI at: http://localhost:8080/swagger/index.html"
```

## Complete Example

### API Handler with Swagger Annotations

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler is a sample API handler
type UserHandler struct{}

// User represents a user in the system
type User struct {
	ID       string `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Age      int    `json:"age" example:"25"`
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Routes registers the routes for the user handler
func (h *UserHandler) Routes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("", h.ListUsers)
		userGroup.GET("/:id", h.GetUser)
		userGroup.POST("", h.CreateUser)
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.DELETE("/:id", h.DeleteUser)
	}
}

// ListUsers godoc
// @Summary      List users
// @Description  get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   User
// @Failure      500  {object}  map[string]interface{}
// @Router       /users [get]
// @Security     BearerAuth
func (h *UserHandler) ListUsers(c *gin.Context) {
	users := []User{
		{ID: "1", Username: "johndoe", Email: "john@example.com", Age: 25},
		{ID: "2", Username: "janedoe", Email: "jane@example.com", Age: 30},
	}
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get a user
// @Description  get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  User
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [get]
// @Security     BearerAuth
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user := User{ID: id, Username: "johndoe", Email: "john@example.com", Age: 25}
	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create a user
// @Description  create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User object"
// @Success      201   {object}  User
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users [post]
// @Security     BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = "3" // In a real app, this would be generated
	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  update an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "User ID"
// @Param        user  body      User    true  "User object"
// @Success      200   {object}  User
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/{id} [put]
// @Security     BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = id
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  delete an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [delete]
// @Security     BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
```

### Main Application with Swagger

```go
package main

import (
	"context"
	"log"

	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	_ "path/to/your/docs" // Import docs for swagger
)

// HealthChecker is a simple health checker
type HealthChecker struct{}

// Health implements the Healthy interface
func (h *HealthChecker) Health(ctx context.Context) error {
	return nil
}

func main() {
	// Create a new user handler
	userHandler := NewUserHandler()

	// Configure Swagger
	err := sfrouting.RegisterConnection(
		sfrouting.WithSwagger(sfrouting.SwaggerConfig{
			Enabled:  true,
			Title:    "SF Routing API",
			Version:  "1.0",
			Host:     "localhost:8080",
			BasePath: "/api/v1",
			Schemes:  []string{"http", "https"},
		}),
		sfrouting.WithHealthChecks(&HealthChecker{}),
	)
	if err != nil {
		log.Fatalf("Failed to register connection: %v", err)
	}

	// Register API routes
	sfrouting.RegisterRouterGroup("/api/v1", userHandler)

	// Start the server
	log.Println("Starting server on :8080")
	if err := sfrouting.StartServer(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

## Usage

### Complete Example

```go
package main

import (
	"context"
	"net/http"
	"time"

	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	"git.snappfood.ir/backend/go/packages/sf-routing/middlewares"
	"github.com/gin-gonic/gin"
)

// MyHealthCheck implements the Healthy interface
type MyHealthCheck struct{}

func (h *MyHealthCheck) Health(ctx context.Context) error {
	// Simulate a health check
	time.Sleep(100 * time.Millisecond)
	return nil
}

// MyPrometheusExporter implements the PrometheusExporter interface
type MyPrometheusExporter struct{}

func (p *MyPrometheusExporter) Handler() http.Handler {
	// In a real implementation, you would use the Prometheus client library 
	// For example: promhttp.Handler() from github.com/prometheus/client_golang/prometheus/promhttp
	
	// This is a simple example that returns a basic HTTP handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`# HELP example_metric Example metric
# TYPE example_metric gauge
example_metric 42`))
	})
}

// MyLogger implements the Logger interface
type MyLogger struct{}

func (l *MyLogger) Debug(msg string, extra map[string]interface{})                                { /* Implementation */ }
func (l *MyLogger) Info(msg string, extra map[string]interface{})                                 { /* Implementation */ }
func (l *MyLogger) Warn(msg string, extra map[string]interface{})                                 { /* Implementation */ }
func (l *MyLogger) Error(msg string, extra map[string]interface{})                                { /* Implementation */ }
func (l *MyLogger) Fatal(msg string, extra map[string]interface{})                                { /* Implementation */ }
func (l *MyLogger) Debugf(template string, args ...interface{})                                   { /* Implementation */ }
func (l *MyLogger) Infof(template string, args ...interface{})                                    { /* Implementation */ }
func (l *MyLogger) Warnf(template string, args ...interface{})                                    { /* Implementation */ }
func (l *MyLogger) Errorf(template string, args ...interface{})                                   { /* Implementation */ }
func (l *MyLogger) Fatalf(template string, args ...interface{})                                   { /* Implementation */ }
func (l *MyLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{})    { /* Implementation */ }
func (l *MyLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{})     { /* Implementation */ }
func (l *MyLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{})     { /* Implementation */ }
func (l *MyLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{})    { /* Implementation */ }
func (l *MyLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{})    { /* Implementation */ }
func (l *MyLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) { /* Implementation */ }
func (l *MyLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{})  { /* Implementation */ }
func (l *MyLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{})  { /* Implementation */ }
func (l *MyLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) { /* Implementation */ }
func (l *MyLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) { /* Implementation */ }

// HomeController handles home routes
type HomeController struct{}

func (ctrl *HomeController) Routes(router *gin.Engine) {
	// Using Mix function to combine handler with middleware
	router.GET("/", middlewares.Mix(ctrl.Home, middlewares.LoggedinMiddle, middlewares.LightMode))
	router.GET("/mobile/v3/user/new-home", middlewares.Mix(ctrl.NewHome, middlewares.LoggedinMiddle, middlewares.CheckDisasterMiddle, middlewares.LightMode))
	router.GET("/search/api/v5/home", middlewares.Mix(ctrl.SearchHome, middlewares.LoggedinMiddle, middlewares.LightMode, middlewares.ABTest))
}

func (ctrl *HomeController) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the home page"})
}

func (ctrl *HomeController) NewHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the new home page"})
}

func (ctrl *HomeController) SearchHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the search home page"})
}

// UserController handles user routes
type UserController struct{}

func (ctrl *UserController) Routes(router *gin.RouterGroup) {
	// Using router group for user routes
	router.GET("/profile", middlewares.Mix(ctrl.Profile, middlewares.LoggedinMiddle))
	router.GET("/settings", middlewares.Mix(ctrl.Settings, middlewares.LoggedinMiddle))
	router.POST("/update", middlewares.Mix(ctrl.Update, middlewares.LoggedinMiddle))
}

func (ctrl *UserController) Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User profile"})
}

func (ctrl *UserController) Settings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User settings"})
}

func (ctrl *UserController) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// RegisterRoutes registers all routes
func RegisterRoutes() {
	// Register main routes
	sfrouting.RegisterRouter(&HomeController{})

	// Register user routes in a group
	sfrouting.RegisterRouterGroup("/api/user", &UserController{})
}

func main() {
	// Configure the server
	sfrouting.RegisterConnection(
		sfrouting.WithLogger(&MyLogger{}),
		sfrouting.WithHealthChecks(&MyHealthCheck{}),
		sfrouting.WithPrometheusExporter(&MyPrometheusExporter{}, sfrouting.PrometheusConfig{
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

	// Register routes
	RegisterRoutes()

	// Start the server
	sfrouting.StartServer(":8080")
}

## Logger Usage

This package provides a comprehensive logging interface that can be used throughout your application. Here's how to use it:

### Setting up the logger

```go
package main

import (
	"github.com/yourlogger/implementation"
	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
)

func main() {
	// Create your logger implementation that satisfies the Logger interface
	myLogger := implementation.NewLogger()
	
	// Pass it to the registry using WithLogger
	err := sfrouting.RegisterConnection(
		sfrouting.WithLogger(myLogger),
		// Other options...
	)
	
	if err != nil {
		panic("Failed to register connection: " + err.Error())
	}
	
	// Start your service
	sfrouting.StartServer(":8080")
}
```

### Using the logger in your code

```go
package mypackage

import (
	"context"
	"git.snappfood.ir/backend/go/packages/sf-routing"
)

func SomeFunction(logger sfrouting.Logger) {
	// Log at different levels
	logger.Debug("Debug message", map[string]interface{}{"key": "value"})
	logger.Info("Info message", map[string]interface{}{"key": "value"})
	logger.Warn("Warning message", map[string]interface{}{"key": "value"})
	logger.Error("Error message", map[string]interface{}{"key": "value"})
	
	// Using categories and subcategories with proper typing
	extraMap := map[string]interface{}{
		sfrouting.ExtraKey.Network.HostIP: "127.0.0.1",
	}
	
	// Use the predefined categories and subcategories directly
	logger.InfoWithCategory(
		sfrouting.Category.System.General,
		sfrouting.SubCategory.Operation.Startup,
		"Application starting",
		extraMap,
	)
	
	// Using context-aware logging
	ctx := context.Background()
	logger.InfoContext(ctx, "Context-aware info", map[string]interface{}{"key": "value"})
	
	// Using formatted logging
	logger.Infof("Hello %s", "world")
}
```

### Always check for nil logger

When using the logger from the registry, always check that it's not nil:

```go
func SomeFunction() {
	registry := // get registry
	
	if registry.logger != nil {
		extraMap := map[string]interface{}{
			sfrouting.ExtraKey.HTTP.Path: "/api/users",
		}
		
		registry.logger.InfoWithCategory(
			sfrouting.Category.System.General,
			sfrouting.SubCategory.Operation.Startup, 
			"Message",
			extraMap,
		)
	}
}
```

### Example of how to use the logger with structured categories, subcategories, and extra keys

```go
package main

import (
	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	"context"
)

func main() {
	logger := // get logger
	
	// Use the predefined structured categories, subcategories, and extra keys
	extraMap := map[string]interface{}{
		sfrouting.ExtraKey.Request.UserID:      "12345",
		sfrouting.ExtraKey.Request.RequestID:   "req-678",
		sfrouting.ExtraKey.HTTP.StatusCode:     200,
		sfrouting.ExtraKey.Performance.Latency: "50ms",
	}
	
	// Log with structured category and subcategory
	logger.InfoWithCategory(
		sfrouting.Category.API.API,
		sfrouting.SubCategory.API.Response,
		"API response sent",
		extraMap,
	)
	
	// Log an error with appropriate categories
	errorMap := map[string]interface{}{
		sfrouting.ExtraKey.Error.ErrorMessage: "Connection timeout",
		sfrouting.ExtraKey.Error.ErrorCode:    "E1001",
	}
	
	logger.ErrorWithCategory(
		sfrouting.Category.Error.Error,
		sfrouting.SubCategory.Status.Timeout,
		"Failed to connect to database",
		errorMap,
	)
}
```

## Feature Details

### Health Check Registration

```go
// MyHealthCheck implements the Healthy interface
type MyHealthCheck struct{}

func (h *MyHealthCheck) Health(ctx context.Context) error {
	// Implement health check logic
	return nil
}

// Register health checks
sfrouting.RegisterHealthCheck(&MyHealthCheck{})
```

### Route Registration with Middleware

```go
// HomeController handles home routes
type HomeController struct{}

func (ctrl *HomeController) Routes(router *gin.Engine) {
	router.GET("/", middlewares.Mix(ctrl.Home, middlewares.LoggedinMiddle, middlewares.LightMode))
}

func (ctrl *HomeController) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the home page"})
}

// Register routes
sfrouting.RegisterRouter(&HomeController{})
```

### Router Groups

```go
// UserController handles user routes
type UserController struct{}

func (ctrl *UserController) Routes(router *gin.RouterGroup) {
	router.GET("/profile", middlewares.Mix(ctrl.Profile, middlewares.LoggedinMiddle))
	router.GET("/settings", middlewares.Mix(ctrl.Settings, middlewares.LoggedinMiddle))
}

func (ctrl *UserController) Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User profile"})
}

func (ctrl *UserController) Settings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User settings"})
}

// Register routes in a group
sfrouting.RegisterRouterGroup("/api/user", &UserController{})
```

### Global Middleware

```go
// Configure the server with global middleware
sfrouting.RegisterConnection(
	sfrouting.WithGlobalMiddleware(middlewares.LoggedinMiddle),
)
```

### Custom Error Handler

```go
// Custom error handler
func MyErrorHandler(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

// Configure the server with custom error handler
sfrouting.RegisterConnection(
	sfrouting.WithErrorHandler(MyErrorHandler),
)
```

### Swagger Documentation

```go
// Configure the server with Swagger
sfrouting.RegisterConnection(
	sfrouting.WithSwagger(sfrouting.SwaggerConfig{
		Enabled:  true,
		Title:    "My API",
		Version:  "1.0",
		Host:     "api.example.com",
		BasePath: "/v1",
		Path:     "/docs/*any",
		Schemes:  []string{"https"},
	}),
)
```

You can also use the default configuration and modify only the settings you need:

```go
// Get default Swagger configuration
config := sfrouting.DefaultSwaggerConfig()

// Modify only the settings you need
config.Enabled = true
config.Title = "My API"
config.Version = "1.0"

// Configure the server with modified Swagger settings
sfrouting.RegisterConnection(
	sfrouting.WithSwagger(config),
)
```

### Custom Logger

```go
// MyLogger implements the Logger interface
type MyLogger struct{}

func (l *MyLogger) Debug(msg string, extra map[string]interface{})                                { /* Implementation */ }
func (l *MyLogger) Info(msg string, extra map[string]interface{})                                 { /* Implementation */ }
func (l *MyLogger) Warn(msg string, extra map[string]interface{})                                 { /* Implementation */ }
func (l *MyLogger) Error(msg string, extra map[string]interface{})                                { /* Implementation */ }
func (l *MyLogger) Fatal(msg string, extra map[string]interface{})                                { /* Implementation */ }
func (l *MyLogger) Debugf(template string, args ...interface{})                                   { /* Implementation */ }
func (l *MyLogger) Infof(template string, args ...interface{})                                    { /* Implementation */ }
func (l *MyLogger) Warnf(template string, args ...interface{})                                    { /* Implementation */ }
func (l *MyLogger) Errorf(template string, args ...interface{})                                   { /* Implementation */ }
func (l *MyLogger) Fatalf(template string, args ...interface{})                                   { /* Implementation */ }
func (l *MyLogger) DebugContext(ctx context.Context, msg string, extra map[string]interface{})    { /* Implementation */ }
func (l *MyLogger) InfoContext(ctx context.Context, msg string, extra map[string]interface{})     { /* Implementation */ }
func (l *MyLogger) WarnContext(ctx context.Context, msg string, extra map[string]interface{})     { /* Implementation */ }
func (l *MyLogger) ErrorContext(ctx context.Context, msg string, extra map[string]interface{})    { /* Implementation */ }
func (l *MyLogger) FatalContext(ctx context.Context, msg string, extra map[string]interface{})    { /* Implementation */ }
func (l *MyLogger) DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{}) { /* Implementation */ }
func (l *MyLogger) InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{})  { /* Implementation */ }
func (l *MyLogger) WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{})  { /* Implementation */ }
func (l *MyLogger) ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{}) { /* Implementation */ }
func (l *MyLogger) FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{}) { /* Implementation */ }

// Example of how to use the logger with structured categories, subcategories, and extra keys
func ExampleLoggerUsage() {
	logger := &MyLogger{}
	
	// Create a map with the predefined extra keys
	extraMap := map[string]interface{}{
		sfrouting.ExtraKey.Network.HostIP:   "localhost",
		sfrouting.ExtraKey.HTTP.ContentType: "application/json",
	}
	
	// Log with proper category and subcategory constants
	logger.InfoWithCategory(
		sfrouting.Category.System.General,
		sfrouting.SubCategory.Operation.Startup,
		"Application is starting",
		extraMap,
	)
}

// Configure the server with custom logger
sfrouting.RegisterConnection(
	sfrouting.WithLogger(&MyLogger{}),
)
```

### Direct Gin Configuration

You can directly configure the Gin engine using the `WithGinConfig` option:

```go
// Configure the server with direct Gin configuration
sfrouting.RegisterConnection(
	sfrouting.WithGinConfig(func(engine *gin.Engine) {
		// Set Gin to release mode
		gin.SetMode(gin.ReleaseMode)
		
		// Configure custom recovery middleware
		engine.Use(gin.Recovery())
		
		// Set custom trusted proxies
		engine.SetTrustedProxies([]string{"192.168.1.1"})
		
		// Configure custom HTML templates
		engine.LoadHTMLGlob("templates/*")
		
		// Configure custom static files
		engine.Static("/static", "./static")
		
		// Any other Gin configuration
	}),
)
```

### CORS Configuration

You can configure CORS (Cross-Origin Resource Sharing) settings using the `WithCorsConfig` option:

```go
// Configure the server with CORS
sfrouting.RegisterConnection(
	sfrouting.WithCorsConfig(sfrouting.CorsConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://example.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}),
)
```

### Predefined Middleware

The library includes several predefined middleware functions that you can use in your routes:

```go
// LightMode middleware for light mode
middlewares.LightMode

// ABTest middleware for A/B testing
middlewares.ABTest

// CheckDisasterMiddle middleware for disaster checking
middlewares.CheckDisasterMiddle

// LoggedinMiddle middleware for checking if user is logged in
middlewares.LoggedinMiddle

// Mix combines a handler with middleware
middlewares.Mix(handler, middleware1, middleware2, ...)
```

You can also create your own middleware functions:

```go
// Custom middleware
func MyCustomMiddleware(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Do something before the handler
		c.Set("custom", "value")
		handler(c)
		// Do something after the handler
	}
}

// Use the custom middleware
router.GET("/custom", middlewares.Mix(ctrl.Custom, MyCustomMiddleware))
```

### Creating custom categories and keys

The logger package now provides structured constants for categories, subcategories, and extra keys:

```go
package main

import (
	sfrouting "git.snappfood.ir/backend/go/packages/sf-routing"
	"context"
)

func main() {
	logger := // get logger
	
	// Use the predefined structured categories, subcategories, and extra keys
	extraMap := map[string]interface{}{
		sfrouting.ExtraKey.Request.UserID:      "12345",
		sfrouting.ExtraKey.Request.RequestID:   "req-678",
		sfrouting.ExtraKey.HTTP.StatusCode:     200,
		sfrouting.ExtraKey.Performance.Latency: "50ms",
	}
	
	// Log with structured category and subcategory
	logger.InfoWithCategory(
		sfrouting.Category.API.API,
		sfrouting.SubCategory.API.Response,
		"API response sent",
		extraMap,
	)
	
	// Log an error with appropriate categories
	errorMap := map[string]interface{}{
		sfrouting.ExtraKey.Error.ErrorMessage: "Connection timeout",
		sfrouting.ExtraKey.Error.ErrorCode:    "E1001",
	}
	
	logger.ErrorWithCategory(
		sfrouting.Category.Error.Error,
		sfrouting.SubCategory.Status.Timeout,
		"Failed to connect to database",
		errorMap,
	)
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

### Health Check Registration with Option

You can now register health checks in two ways:

```go
// Old way (still supported for backward compatibility)
sfrouting.RegisterHealthCheck(&MyHealthCheck{})

// New way through RegisterConnection
sfrouting.RegisterConnection(
	sfrouting.WithHealthChecks(&MyHealthCheck{}, &AnotherHealthCheck{}),
	// Other options...
)
```

### Prometheus Metrics Exporter

To expose Prometheus metrics for your application, implement the `PrometheusExporter` interface and register it using `WithPrometheusExporter`:

```go
// PrometheusExporter defines the interface for Prometheus exporter
type PrometheusExporter interface {
	// Handler returns an http.Handler that will handle Prometheus metrics requests
	Handler() http.Handler
}

// Implement your custom Prometheus exporter
type MyPrometheusExporter struct{}

func (p *MyPrometheusExporter) Handler() http.Handler {
	// In a real implementation, you would use the Prometheus client library
	// For example, with github.com/prometheus/client_golang/prometheus:
	
	// Create a registry
	registry := prometheus.NewRegistry()
	
	// Register your metrics with the registry
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_counter",
		Help: "This is my counter",
	})
	registry.MustRegister(counter)
	
	// Increment counter in your application code
	counter.Inc()
	
	// Return the HTTP handler that serves the metrics
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	
	// Or for a simple implementation:
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`# HELP example_metric Example metric
# TYPE example_metric gauge
example_metric 42`))
	})
}

// Register the exporter with the registry
sfrouting.RegisterConnection(
	sfrouting.WithPrometheusExporter(&MyPrometheusExporter{}, sfrouting.PrometheusConfig{
		Enabled: true,
		Path:    "/metrics",
	}),
	// Other options...
)
```

The library automatically creates a route at the specified path (default: `/metrics`) that uses your HTTP handler to serve Prometheus metrics.

You can also use the default configuration and modify only the settings you need:

```go
// Get default Prometheus configuration
config := sfrouting.DefaultPrometheusConfig()

// Modify only the settings you need
config.Enabled = true
config.Path = "/custom-metrics"

// Configure the server with modified Prometheus settings
sfrouting.RegisterConnection(
	sfrouting.WithPrometheusExporter(&MyPrometheusExporter{}, config),
	// Other options...
)
```