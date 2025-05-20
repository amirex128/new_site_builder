# SF-Redis

A resilient Redis connection registry for Go.

[![Go Reference](https://pkg.go.dev/badge/git.snappfood.ir/backend/go/packages/sf-redis.svg)](https://pkg.go.dev/git.snappfood.ir/backend/go/packages/sf-redis)
[![Go Report Card](https://goreportcard.com/badge/git.snappfood.ir/backend/go/packages/sf-redis)](https://goreportcard.com/report/git.snappfood.ir/backend/go/packages/sf-redis)

This library provides a registry for managing multiple Redis connections in Go applications:

- **Resilient Connection Management**: Automatic retry and reconnection with exponential backoff, panic recovery, and structured logging
- **Multiple Redis Support**: Connect to multiple Redis instances
- **Service Registry Integration**: Works with SF Service Registry
- **Health Checks**: Comprehensive health check functionality
- **Connection Pooling**: Configurable connection pools

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Retry Options & Resilience](#retry-options--resilience)
- [Connection Registry](#connection-registry)
- [Service Connector](#service-connector)
- [Health Checks](#health-checks)
- [Features](#features)
- [Redis Operations](#redis-operations)
- [Connection Options](#connection-options)
- [API Reference](#api-reference)
- [Implementation Notes](#implementation-notes)

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-redis
```

## Quick Start

```go
package main

import (
	"context"
	"log"
	"time"

	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

func main() {
	ctx := context.Background()

	// Register SfRedis connections with retry options for resilience
	err := sfredis.RegisterConnection(
		sfredis.WithGlobalOptions(func(options *sfredis.Options) {
			options.PoolSize = 20
			options.MinIdleConns = 5
			options.MaxRetries = 3
		}),
		sfredis.WithRetryOptions(&sfredis.RetryOptions{
			MaxRetries:     5,
			InitialBackoff: time.Second,
			MaxBackoff:     15 * time.Second,
			BackoffFactor:  1.5,
		}),
		sfredis.WithConnectionDetails("cache", "localhost:6379", "", 0),
		sfredis.WithConnectionDetails("session", "localhost:6380", "password", 1),
	)

	if err != nil {
		log.Fatalf("Failed to register SfRedis connections: %v", err)
	}

	// Get the Redis client by name
	cacheClient := sfredis.MustClient(ctx, "cache")
	sessionClient := sfredis.MustClient(ctx, "session")

	cacheClient.Set(ctx, "key1", "value1", time.Hour)
	value, err := cacheClient.Get(ctx, "key1")
	if err != nil {
		log.Printf("Get failed: %v\n", err)
	}

	// Close all connections when application exits
	defer sfredis.CloseConnections()
}
```

## Retry Options & Resilience

SF-Redis provides a robust retry mechanism for all connection attempts to external data sources. This includes:

- **Exponential Backoff**: Retry intervals increase exponentially up to a maximum.
- **Panic Recovery**: All panics during connection attempts are caught and logged as errors, preventing application crashes.
- **Structured Logging**: All retry attempts, successes, failures, and panics are logged with categories and extra context.

You can configure retry behavior globally for all connections using `WithRetryOptions`:

```go
sfredis.WithRetryOptions(&sfredis.RetryOptions{
	MaxRetries:     5,                // Maximum number of retry attempts
	InitialBackoff: time.Second,      // Initial waiting time between retries
	MaxBackoff:     15 * time.Second, // Maximum waiting time between retries
	BackoffFactor:  1.5,              // Exponential backoff multiplier
})
```

**Example:**

```go
err := sfredis.RegisterConnection(
	sfredis.WithRetryOptions(&sfredis.RetryOptions{
		MaxRetries:     5,
		InitialBackoff: time.Second,
		MaxBackoff:     15 * time.Second,
		BackoffFactor:  1.5,
	}),
	// ... other options ...
)
```

## Connection Registry
The connection registry is the core of SF-Redis. It allows you to register Redis connections by name in one place (typically in main.go), then reuse them throughout your application without creating new connections each time.

## Service Connector

SF-Redis provides a ServiceConnector that implements the serviceregistry.ServiceConnector interface, making it easy to integrate with the SF Service Registry. You can also pass retry options through the connector:

```go
package main

import (
	"context"
	"log"
	"time"

	"git.snappfood.ir/backend/go/packages/sf-service-registry"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

func main() {
	ctx := context.Background()
	connector := sfredis.GetConnector()

	err := serviceregistry.RegisterServices(
		nil,       // No MySQL
		connector, // Redis
		nil,       // No Elastic
		nil,       // No RabbitMQ
		nil,       // No SQL Server
	)
	if err != nil {
		log.Fatalf("Failed to register services: %v", err)
	}

	// Register connections with retry options
	err = connector.RegisterConnection(
		sfredis.WithRetryOptions(&sfredis.RetryOptions{
			MaxRetries:     5,
			InitialBackoff: time.Second,
			MaxBackoff:     15 * time.Second,
			BackoffFactor:  1.5,
		}),
		// ... other options ...
	)
	if err != nil {
		log.Fatalf("Failed to register SfRedis connections: %v", err)
	}

	cacheClient := sfredis.MustClient(ctx, "cache")
	// Use the connection...
}
```

## Health Checks

SF-Redis provides a Health function that checks if all Redis connections are alive:

```go
package main

import (
    "context"
    "log"
    "time"
    
    sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

func main() {
    // Register your connections...
    
    // Check health of all connections
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err := sfredis.Health(ctx)
    if err != nil {
        log.Fatalf("Health check failed: %v", err)
    }
    
    log.Println("All Redis connections are healthy")
}
```

### Health Checks with Gin Framework

You can integrate health checks with a Gin web server:

```go
package main

import (
    "context"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"

	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

func main() {
    // Register your connections...
    
    // Set up Gin router
    r := gin.Default()
    
    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        // Create a context with timeout for health checks
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        // Check Redis connections
        if err := sfredis.Health(ctx); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unhealthy",
                "error": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
        })
    })
    
    // Start server
    r.Run(":8080")
}
```

## Features

- **Named Connection Registry** - Register multiple Redis connections by name
- **Lazy Connection Initialization** - Connections are only established when first used
- **Connection Pool Configuration** - Configure idle connections, max open connections, etc.
- **Connection-Specific Options** - Set custom options for each connection
- **Global Options** - Apply settings to all connections
- **Thread-Safe Implementation** - Proper locking for concurrent access
- **Connection Management** - Close and reopen connections as needed
- **Health Checks** - Comprehensive health check functionality
- **Service Registry Integration** - Works with SF Service Registry
- **Application Performance Monitoring (APM) integration**
- **Resilient Retry System** - Exponential backoff, panic recovery, and structured logging for all connection attempts

## SfRedis custom Operations

When you have registered your connections, you can use them anywhere in your application:

```go
package service

import (
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
)

```

## Connection Options

```go
// Configure Redis logger
customLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold:             200 * time.Millisecond,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: false,
        Colorful:                  true,
    },
)

sf_redis.RegisterConnection(
    sf_redis.WithLogger(customLogger),
    sf_redis.WithConnectionDetails("cache", "localhost:6379", "", 0),
)
```

## API Reference

### Connection Registration and Management

- `RegisterConnection(opts ...RegistryOption) error`
- `MustClient(ctx context.Context, name string) *SfRedis` (returns nil and logs error if connection fails)
- `SafeClient(ctx context.Context, name string) (*SfRedis, error)`
- `CloseConnection(name string) error`
- `CloseConnections()`
- `Health(ctx context.Context) error`

### Service Connector

- `GetConnector() *ServiceConnector`
- `ServiceConnector.WithOptions(opts ...RegistryOption) *ServiceConnector`
- `ServiceConnector.Health(ctx context.Context) error`
- `ServiceConnector.RegisterConnection(opts ...interface{}) error` (supports WithRetryOptions)

### Options

- `WithLogger(logger Logger) RegistryOption`
- `WithGlobalOptions(options ...RedisOption) RegistryOption`
- `WithConnectionDetails(name, addr, password string, db int, options ...RedisOption) RegistryOption`
- `WithRetryOptions(options *RetryOptions) RegistryOption`

#### RetryOptions struct

```go
type RetryOptions struct {
	MaxRetries     int           // Maximum number of retry attempts
	InitialBackoff time.Duration // Starting backoff duration
	MaxBackoff     time.Duration // Maximum backoff duration
	BackoffFactor  float64       // Multiplier for each subsequent retry
}
```

## Implementation Notes

- The retry system is implemented with exponential backoff and panic recovery for all connection attempts.
- All panics during connection attempts are caught and logged as errors, preventing application crashes.
- All retry attempts, successes, failures, and panics are logged using the provided logger and categories.
- The library uses the `github.com/go-redis/redis/v8` package for Redis operations.

### Configuration

Services are configured in `services.json`:

```json
{
  "redis": [
    {
      "name": "cache",
      "env_host": "host",
      "env_port": "port",
      "env_password": "password",
      "env_database": "db",
      "env_pool_size": "pool_size",
      "env_min_idle_conns": "min_idle_conns"
    }
  ]
}
```

The system uses Viper to read environment variables based on the configuration. For example, for a Redis service named `cache`, it will look for environment variables like:
- `cache_host`
- `cache_port`
- `cache_password`
- `cache_database`
- `cache_pool_size`
- `cache_min_idle_conns`

## APM Integration

SfRedis includes built-in support for Application Performance Monitoring (APM) through Elastic APM. To enable APM monitoring, you need to set the following environment variables:

```bash
# Required environment variables for APM
ELASTIC_APM_SERVICE_NAME=your-service-name
ELASTIC_APM_SERVER_URL=http://your-apm-server:8200
ELASTIC_APM_SECRET_TOKEN=your-secret-token
ELASTIC_APM_ENVIRONMENT=production
```

The APM integration provides:
- Automatic tracing of Redis operations
- Performance metrics for Redis commands
- Error tracking and reporting
- Distributed tracing support