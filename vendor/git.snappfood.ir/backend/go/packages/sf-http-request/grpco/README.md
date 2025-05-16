# SF gRPC Client (grpco)

A powerful, fluent gRPC client for Go applications with built-in connection management, resilience patterns, and observability.

## Features

- Connection registry for managing gRPC connections
- Fluent API for making gRPC requests
- Resilience patterns: timeout, retry, fallback
- Automatic connection retries
- Observability with structured logging
- Header management with global and connection-specific headers
- Health checks for monitoring connection status

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-http-request
```

## Quick Start

### Register Connections and Services

```go
import (
    "google.golang.org/grpc/metadata"
    "git.snappfood.ir/backend/go/packages/sf-http-request/grpco"
)

func init() {
    // Define services for the connection
    svcDef := map[string]grpco.ServiceDefinition{
        "story": {
            ClientConstructor: service.NewStoryServiceClient,
            Methods: map[string]string{
                "List": "/service.StoryService/List",
                "Get": "/service.StoryService/Get",
            },
        },
    }

    // Register a connection with services in a single call
    grpco.RegisterConnection(
        // Define the connection details and pass services directly
        grpco.WithConnectionDetails(
            "story-service", 
            "story-service.example.com:50051",
            grpco.WithInsecure(),
            svcDef, // Pass services as an argument 
        ),
        
        // Add global headers for all connections
        grpco.WithGlobalHeaders(metadata.MD{
            "x-app-name": []string{"my-service"},
            "x-version": []string{"1.0"},
        }),
        
        // Set logger for observability
        grpco.WithLogger(logger),
    )
    
    // Alternatively, you can still register services separately
    // using WithServices if preferred
    otherSvcDef := map[string]grpco.ServiceDefinition{
        "user": {
            ClientConstructor: service.NewUserServiceClient,
            Methods: map[string]string{
                "GetUser": "/service.UserService/GetUser",
            },
        },
    }
    
    grpco.RegisterConnection(
        grpco.WithServices("user-service", otherSvcDef),
    )
}
```

### Make gRPC Requests with Fluent Interface

```go
// Create a request and get a response
req := &service.ListRequest{
    UserLatitude:  35.7219,
    UserLongitude: 51.3347,
}
var resp service.ListResponse

// Make the request with resilience patterns
// Use FromService with 3 parameters: connectionName, serviceName, methodName
err := grpco.FromService("story-service", "story", "List")
    .Request(req)
    .Header("x-correlation-id", "abc-123")
    .Timeout(5000) // 5 seconds
    .Retry(3, 500) // Retry 3 times with 500ms delay
    .Fallback(func() {
        fmt.Println("Using fallback data")
    })
    .Send(&resp)

// Use response
fmt.Printf("Received %d reels\n", len(resp.Reels))
```

## Detailed Usage

### Register Multiple Connections

```go
// Create service definitions for connections
storyServices := map[string]grpco.ServiceDefinition{
    "story": {
        ClientConstructor: service.NewStoryServiceClient,
        Methods: map[string]string{
            "List": "/service.StoryService/List",
            "Get":  "/service.StoryService/Get",
        },
    },
}

userServices := map[string]grpco.ServiceDefinition{
    "user": {
        ClientConstructor: service.NewUserServiceClient,
        Methods: map[string]string{
            "GetUser":    "/service.UserService/GetUser",
            "ListUsers":  "/service.UserService/ListUsers",
        },
    },
}

// Register connections with services directly
grpco.RegisterConnection(
    grpco.WithConnectionDetails(
        "story-service", 
        "story-service:50051",
        grpco.WithInsecure(),
        storyServices, // Pass services directly
    ),
)

grpco.RegisterConnection(
    grpco.WithConnectionDetails(
        "user-service", 
        "user-service:50052",
        grpco.WithTLS(&tls.Config{...}),
        userServices, // Pass services directly
    ),
)

// You can also register multiple service maps for a single connection
authServices := map[string]grpco.ServiceDefinition{
    "auth": {
        ClientConstructor: service.NewAuthServiceClient,
        Methods: map[string]string{
            "Login":  "/service.AuthService/Login",
            "Logout": "/service.AuthService/Logout",
        },
    },
}

// Register multiple service maps at once
grpco.RegisterConnection(
    grpco.WithConnectionDetails(
        "multi-service", 
        "multi-service:50053",
        grpco.WithInsecure(),
        storyServices,  // First service map
        userServices,   // Second service map
        authServices,   // Third service map
    ),
)
```

### Connection-Specific Headers

```go
// Set headers for a specific connection
grpco.SetConnectionHeader("story-service", "x-version", "1.0")

// Set multiple headers
grpco.SetConnectionHeaders("story-service", map[string]string{
    "x-version": "1.0",
    "x-region": "us-east-1",
})
```

### Making gRPC Requests

```go
// Basic request pattern - no need to check error after FromService
var listResp service.ListResponse
grpco.FromService("story-service", "story", "List")
    .Request(&service.ListRequest{UserSessionID: "123"})
    .Send(&listResp) // Error handling happens here

// With resilience options
var getResp service.GetResponse
err := grpco.FromService("story-service", "story", "Get")
    .Request(&service.GetRequest{Id: "456"})
    .Header("x-request-id", "req-789")
    .SetContext(ctx)
    .Timeout(3000)
    .Retry(2, 1000)
    .Fallback(func() {
        // Fallback logic when all retries fail
    })
    .Send(&getResp)
if err != nil {
    // Handle errors from initialization or execution
}

// Using RequestWithResponse for better type safety
getReq := &service.GetRequest{Id: "456"}
var getResponse service.GetResponse
err := grpco.FromService("story-service", "story", "Get")
    .RequestWithResponse(getReq, &getResponse)
    .Timeout(5000)
    .Send()

// For backward compatibility with older code
err := grpco.FromServiceWithKey("story-service", "List")
    .Request(&service.ListRequest{UserSessionID: "123"})
    .Send(&listResp)
```

### Health Checking

```go
// Check health of all connections
err := grpco.Health(context.Background())
if err != nil {
    log.Printf("gRPC health check failed: %v", err)
}
```

### Managing Headers

```go
// Set global headers for all connections
grpco.SetGlobalHeader("x-api-key", "global-api-key")
grpco.SetGlobalHeaders(map[string]string{
    "x-client-id": "client-123",
    "x-source": "backend-service",
})

// Get all global headers
headers := grpco.GetGlobalHeaders()

// Remove a specific global header
grpco.RemoveGlobalHeader("x-api-key")

// Clear all global headers
grpco.ClearGlobalHeaders()
```

### Cleanup

```go
// Close all connections when shutting down
func shutdown() {
    if err := grpco.CloseConnections(); err != nil {
        log.Printf("Error closing gRPC connections: %v", err)
    }
    
    // Or close a specific connection
    if err := grpco.CloseConnection("story-service"); err != nil {
        log.Printf("Error closing story service connection: %v", err)
    }
}
```

## Core API Methods

The library has been simplified to provide a single, clear way to make gRPC requests:

| Method | Description |
|--------|-------------|
| `WithConnectionDetails(name, target, ...interface{})` | Sets connection details and optionally registers services |
| `FromService(connectionName, serviceName, methodName)` | Creates a request builder for a specific service method |
| `FromServiceWithKey(connectionName, methodKey)` | Backward compatibility method for the old two-parameter API |
| `Request().Send()` | Fluent API for configuring and sending requests |

## Configuration Options

| Option | Description |
|--------|-------------|
| `WithConnectionDetails` | Registers a named connection with a target address, options, and services |
| `WithServices` | Registers services for a connection separately |
| `WithDefaultTransport` | Sets the default transport for all connections |
| `WithGlobalHeaders` | Sets headers to be included in all requests |
| `WithLogger` | Sets the logger for observability |
| `WithInsecure` | Configures the connection to use plaintext |
| `WithTLS` | Configures the connection to use TLS |

## License

Proprietary - SnappFood 