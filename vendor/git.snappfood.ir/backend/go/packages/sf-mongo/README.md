# MongoDB Connection Registry for Go

This package provides a robust, thread-safe registry for managing multiple named MongoDB connections in Go, inspired by production-grade patterns. It supports global and per-connection options, retry logic, and structured logging.

## Features
- Register multiple named MongoDB connections
- Global and per-connection options
- Automatic connection retry with backoff
- Health checks for all connections
- Structured logging integration
- Safe concurrent access
- Easy API for getting and closing connections

## Installation

Add the MongoDB Go driver to your project (already included if you used this package):

```
go get go.mongodb.org/mongo-driver/mongo@v1.12.0
```

## Usage Example

```go
import (
    "context"
    "go.mongodb.org/mongo-driver/mongo/options"
    "git.snappfood.ir/backend/go/packages/sf-mongo"
)

// Implement the Logger interface from logger.go
var logger sfmongo.Logger = ...

func main() {
    // Register a MongoDB connection
    err := sfmongo.RegisterConnection(
        sfmongo.WithLogger(logger),
        sfmongo.WithConnectionDetails(
            "main-db",
            options.Client().ApplyURI("mongodb://localhost:27017"),
        ),
        sfmongo.WithRetryOptions(sfmongo.DefaultRetryOptions()),
    )
    if err != nil {
        panic(err)
    }

    // Get a MongoDB client
    client, err := sfmongo.GetConnection("main-db")
    if err != nil {
        panic(err)
    }
    defer sfmongo.CloseConnection("main-db")

    // Use the client...
    db := client.Database("test")
    // ...

    // Health check
    if err := sfmongo.Health(context.Background()); err != nil {
        logger.Error("MongoDB health check failed", map[string]interface{}{"error": err.Error()})
    }
}
```

## API

### Registering Connections
- `RegisterConnection(opts ...MongoRegistryOption) error` — Register and start connections.
- `WithLogger(logger Logger)` — Set a logger.
- `WithConnectionDetails(name string, config *options.ClientOptions, options ...interface{})` — Register a named connection.
- `WithRetryOptions(options *RetryOptions)` — Set retry options.
- `WithGlobalOptions(option MongoClientOption)` — Set global client options.

### Getting and Using Clients
- `GetConnection(name string) (*mongo.Client, error)` — Get a client by name (waits for connection).
- `MustMongoClient(name string) *mongo.Client` — Get a client or panic.

### Health Checks
- `Health(ctx context.Context) error` — Check all connections.

### Closing Connections
- `CloseConnections()` — Close all connections.
- `CloseConnection(name string) error` — Close a specific connection.

## Logger and Retry
- Integrate your own logger by implementing the `Logger` interface from `logger.go`.
- Use `RetryOptions` from `retry.go` for custom retry logic.

## License
See project root for license information.