# SF HTTP Request

A robust Go HTTP and gRPC client library with built-in resilience patterns including timeouts, retries, and connection management.

## Features

- 🌐 **HTTP and gRPC Clients** - Unified API for both HTTP and gRPC requests
- 🔄 **Automatic Retries** - Configurable retry logic with backoff
- 🔌 **Connection Registry** - Centralized management of service connections
- 🕒 **Timeouts** - Fine-grained request timeout control
- 🛡️ **Fallbacks** - Graceful degradation with fallback functions
- 📊 **Header Management** - Global and per-connection header control
- 🔍 **Response Utilities** - Simplified response handling
- 📝 **Structured Logging** - Comprehensive logging with categories and context

## Installation

```bash
go get git.snappfood.ir/backend/go/packages/sf-http-request
```

## Quick Start

### HTTP Example

```go
package main

import (
    "fmt"
    "log"
    "time"

    "git.snappfood.ir/backend/go/packages/sf-http-request/httpo"
    "git.snappfood.ir/backend/go/packages/sf-http-request/serviceregistry"
)

func main() {
    // Create a logger implementation of serviceregistry.Logger
    logger := setupLogger() // Your custom logger that implements the interface
    
    // Register a connection with logger
    err := httpo.RegisterConnection(
        httpo.WithConnectionDetails("example-api", "https://api.example.com"),
        httpo.WithLogger(logger),
    )
    if err != nil {
        log.Fatalf("Failed to register connection: %v", err)
    }

    // Make a request with resilience patterns
    resp, err := httpo.Get("example-api", "/users")
        .Query("page", "1")
        .Header("Accept", "application/json")
        .Timeout(5000) // 5 seconds
        .Retry(3, 500) // Retry 3 times with 500ms delay
        .Fallback(func() {
            fmt.Println("Using fallback data")
        })
        .Send()

    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }

    // Check response status
    if !resp.IsSuccess() {
        log.Fatalf("Request failed with status: %d", resp.StatusCode)
    }

    // Parse JSON response
    var users []map[string]interface{}
    if err := resp.JSON(&users); err != nil {
        log.Fatalf("Failed to parse JSON: %v", err)
    }

    fmt.Printf("Found %d users\n", len(users))
}
```

### gRPC Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "git.snappfood.ir/backend/go/packages/sf-http-request/grpco"
    "google.golang.org/grpc/metadata"
)

func main() {
    // Define services for the connection
    svcDef := map[string]grpco.ServiceDefinition{
        "story": {
            ClientConstructor: service.NewStoryServiceClient,
            Methods: map[string]string{
                "List": "/service.StoryService/List",
                "Get":  "/service.StoryService/Get",
            },
        },
    }

    // Register a connection with services in a single call
    err := grpco.RegisterConnection(
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
            "x-version":  []string{"1.0"},
        }),
    )
    if err != nil {
        log.Fatalf("Failed to register connection: %v", err)
    }

    // Wait for the connection to be established (timeout after 5 seconds)
    if err := grpco.WaitForConnection("story-service", 5000); err != nil {
        log.Fatalf("Connection not established: %v", err)
    }

    // Create a request and get a response
    req := &service.ListRequest{
        UserLatitude:  35.7219,
        UserLongitude: 51.3347,
    }
    var resp service.ListResponse

    // Make the request with resilience patterns
    // Use FromService with 3 parameters: connectionName, serviceName, methodName
    err = grpco.FromService("story-service", "story", "List")
        .Request(req)
        .Header("x-correlation-id", "abc-123")
        .Timeout(5000) // 5 seconds
        .Retry(3, 500) // Retry 3 times with 500ms delay
        .Fallback(func() {
            fmt.Println("Using fallback data")
        })
        .Send(&resp)
    
    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }

    // Use response
    fmt.Printf("Received %d reels\n", len(resp.Reels))
}
```

## Logging

The library supports structured logging with categories and subcategories:

```go
// First implement the Logger interface
type MyLogger struct {
    // Your fields
}

// Implement all the methods required by the serviceregistry.Logger interface
func (l *MyLogger) InfoWithCategory(cat serviceregistry.Category, subCat serviceregistry.SubCategory, 
    msg string, extra map[serviceregistry.ExtraKey]interface{}) {
    // Your implementation
}

// Then use it when registering connections
logger := &MyLogger{}

// Register HTTP connection with logger
httpo.RegisterConnection(
    httpo.WithConnectionDetails("example-api", "https://api.example.com"),
    httpo.WithLogger(logger),
)

// Register gRPC connection with logger
grpco.RegisterConnection(
    grpco.WithConnectionDetails("user-service", "localhost:50051", grpco.WithInsecure()),
    grpco.WithLogger(logger),
)
```

## HTTP Connection Registry

Manage HTTP connections centrally:

```go
// Register connections
httpo.RegisterConnection(
    httpo.WithConnectionDetails("example-api", "https://api.example.com"),
)

// Register with custom client options
httpo.RegisterConnection(
    httpo.WithConnectionDetails("other-api", "https://api.other.com", 
        httpo.WithTransport(customTransport),
        httpo.WithTimeout(10*time.Second),
    ),
)

// Set global headers for all connections
httpo.SetGlobalHeaders(map[string]string{
    "X-API-Key": "your-api-key",
    "User-Agent": "SF-HTTP-Client/1.0",
})

// Set headers for a specific connection
httpo.SetConnectionHeader("example-api", "Authorization", "Bearer token123")

// Use a registered connection
resp, err := httpo.Get("example-api", "/users").Send()

// Close connections when done
httpo.CloseConnections()
```

## gRPC Connection Registry

Manage gRPC connections centrally:

```go
// Register connections with services in a single call
grpco.RegisterConnection(
    grpco.WithConnectionDetails(
        "story-service", 
        "story-service:50051",
        grpco.WithInsecure(),
        map[string]grpco.ServiceDefinition{
            "story": {
                ClientConstructor: service.NewStoryServiceClient,
                Methods: map[string]string{
                    "List": "/service.StoryService/List",
                    "Get":  "/service.StoryService/Get",
                },
            },
        },
    ),
    grpco.WithGlobalHeaders(metadata.MD{
        "x-api-key": []string{"your-api-key"},
    }),
)

// Enable auto-connection waiting with timeout
grpco.RegisterConnection(
    grpco.WithAutoConnectTimeout(5000), // Wait up to 5 seconds for connections
    // Other options...
)

// Or wait explicitly for a connection
if err := grpco.WaitForConnection("story-service", 5000); err != nil {
    log.Fatalf("Connection not established: %v", err)
}

// Make requests with the fluent API
var resp service.ListResponse
err := grpco.FromService("story-service", "story", "List")
    .Request(&service.ListRequest{})
    .Header("x-correlation-id", "abc-123")
    .Timeout(5000)
    .Send(&resp)

// Close connections when done
grpco.CloseConnections()
```

## Advanced HTTP Features

```go
// POST with JSON body
resp, err := httpo.Post("example-api", "/users")
    .JSONBody(map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    })
    .Send()

// Form data
resp, err := httpo.Post("example-api", "/login")
    .FormData("username", "john")
    .FormData("password", "secret")
    .Send()

// Authentication
resp, err := httpo.Get("example-api", "/profile")
    .BearerAuth("your-token")
    .Send()

// Direct URL (without registry)
resp, err := httpo.URL("GET", "https://api.example.com/users")
    .Query("page", "1")
    .Send()

// Response handling
if resp.IsSuccess() {
    var data map[string]interface{}
    resp.JSON(&data)
    
    // Or get raw response
    bodyStr := resp.MustString()
    contentType := resp.GetContentType()
}
```

## Advanced gRPC Features

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

// Setting metadata
req, _ := grpco.FromService("story-service", "story", "List")
    .Header("x-correlation-id", "abc-123")
    .Header("x-tenant-id", "tenant1")

// Custom context
ctx := context.WithValue(context.Background(), "key", "value")
req.SetContext(ctx)

// Get response headers and trailers
var resp service.ListResponse
req.Send(&resp)
responseHeaders := req.GetResponseHeaders()
responseTrailers := req.GetResponseTrailers()
```

## Health Checks

```go
// Check HTTP connections
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := httpo.Health(ctx)

// Check gRPC connections
err := grpco.Health(ctx)
```

## License

Copyright © SnapFood 