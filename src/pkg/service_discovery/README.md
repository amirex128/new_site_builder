# Service Registry

This package provides a service registry implementation for service discovery. It supports multiple backends:

- Consul
- etcd
- Zookeeper

## Usage

### Configuration

Add the following environment variables to your configuration:

```env
# Service registry configuration
SERVER_SERVICE_REGISTRY_TYPE=consul      # Type of service registry (consul, zookeeper, etcd)
SERVER_SERVICE_REGISTRY_ADDRESS=localhost:8500  # Address of the service registry
SERVER_SERVICE_REGISTRY_TOKEN=my-token   # Authentication token for the service registry (optional)
SERVER_SERVICE_REGISTRY_USERNAME=user    # Username for the service registry (optional)
SERVER_SERVICE_REGISTRY_PASSWORD=pass    # Password for the service registry (optional)
SERVER_SERVICE_NAME=my-service           # Name of this service
SERVER_SERVICE_PORT=8080                 # Port of this service
SERVER_SERVICE_ADDRESS=localhost         # Address of this service
SERVER_SERVICE_TAGS=api,v1               # Tags for this service (comma-separated)
SERVER_SERVICE_HEALTH_CHECK_URL=http://localhost:8080/health  # Health check URL for this service
SERVER_SERVICE_HEALTH_CHECK_TTL=30s      # Health check TTL for this service
```

### Simplified Integration

The easiest way to use the service registry is to use the `InitAndRegisterService` function, which handles all the parsing, initialization, and registration in a single call:

```go
import (
    "go-boilerplate/src/pkg/service_discovery"
    logger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func main() {
    // Initialize your logger
    log := logger.NewLogger()
    
    // Set up service configuration
    serviceConfig := serviceregistry.ServiceRegistryConfig{
        Registry: serviceregistry.RegistryConfig{
            Type:     "consul",
            Address:  "localhost:8500",
            Token:    "my-token",
            Username: "user",
            Password: "pass",
        },
        Service: serviceregistry.ServiceConfig{
            Name:    "my-service",
            Port:    "8080",
            Address: "localhost",
            Tags:    "api,v1",
            Scheme:  "http",
            HealthCheck: serviceregistry.HealthCheck{
                URL:      "http://localhost:8080/health",
                TTL:      "30s",
                Interval: "10s",
            },
        },
        SetupGracefulShutdown: true, // Auto-setup graceful shutdown
    }
    
    // Initialize and register in one step
    serviceID, err := serviceregistry.InitAndRegisterService(log, serviceConfig)
    if err != nil {
        log.Fatal("Failed to initialize service registry", map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    // The service is now registered!
    // If SetupGracefulShutdown was true, the service will automatically
    // deregister on shutdown (Ctrl+C or SIGTERM)
    
    // Start your service...
}
```

This method handles:
- Parsing port from string to int
- Parsing tags from comma-separated string 
- Initializing the registry with the appropriate backend
- Registering your service
- Setting up graceful shutdown (if requested)

### Manual Integration

If you need more control, you can use the individual functions:

```go
import (
    "go-boilerplate/src/pkg/service_discovery"
    logger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func main() {
    log := logger.NewLogger()
    
    // Initialize registry
    config := serviceregistry.Config{
        Address: "localhost:8500",
        Options: map[string]interface{}{
            "token": "my-token",
        },
    }
    
    err := serviceregistry.InitServiceRegistry(log, "consul", config)
    if err != nil {
        // Handle error
    }
    
    // Register service
    serviceInfo := serviceregistry.ServiceInfo{
        Name:          "my-service",
        ID:            serviceregistry.GenerateServiceID("my-service"),
        Address:       "localhost",
        Port:          8080,
        Tags:          []string{"api", "v1"},
        Meta:          map[string]string{"version": "1.0.0"},
        CheckURL:      "http://localhost:8080/health",
        CheckInterval: "10s",
    }
    
    err = serviceregistry.RegisterService(log, serviceInfo)
    if err != nil {
        // Handle error
    }
    
    // Set up deregistration on shutdown
    // ...
}
```

### Discovering Services

To discover a service:

```go
import (
    "fmt"
    "go-boilerplate/src/pkg/service_discovery"
    logger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func CallAnotherService(log logger.Logger) (string, error) {
    // Get URL for the service
    serviceURL, err := serviceregistry.GetServiceURL(log, "another-service")
    if err != nil {
        return "", fmt.Errorf("failed to get service URL: %w", err)
    }

    // Construct the full URL
    fullURL := fmt.Sprintf("%s/api/endpoint", serviceURL)

    // Make the request...
    // ...

    return "Success", nil
}
```

### Making HTTP Requests to Discovered Services

Here's an example of how to make HTTP requests to a discovered service:

```go
import (
    "fmt"
    "net/http"
    "go-boilerplate/src/pkg/service_discovery"
    logger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func CallService(log logger.Logger, serviceName string, path string) (string, error) {
    // Discover the service
    serviceURL, err := serviceregistry.GetServiceURL(log, serviceName)
    if err != nil {
        return "", fmt.Errorf("failed to discover service %s: %w", serviceName, err)
    }

    // Construct the full URL
    fullURL := fmt.Sprintf("%s%s", serviceURL, path)

    // Make a request to the service
    log.Info("Making request to service", map[string]interface{}{
        "service": serviceName,
        "url":     fullURL,
    })

    resp, err := http.Get(fullURL)
    if err != nil {
        return "", fmt.Errorf("failed to call service %s: %w", serviceName, err)
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("service %s returned non-OK status: %d", serviceName, resp.StatusCode)
    }

    return fmt.Sprintf("Successfully called %s service", serviceName), nil
}
```

## Backend-Specific Details

### Consul

Consul is used by default when `RegistryType=consul`. It supports the following features:

- Service registration with health checks
- Service discovery
- Health status monitoring

Example with Consul:

```go
// Initialize service registry with Consul
config := serviceregistry.Config{
    Address: "localhost:8500",
    Options: map[string]interface{}{
        "token": "my-consul-token",
    },
}

// Initialize registry with Consul
err := serviceregistry.InitServiceRegistry(log, "consul", config)
if err != nil {
    return fmt.Errorf("failed to initialize service registry: %w", err)
}

// Register service
serviceInfo := serviceregistry.ServiceInfo{
    Name:          "my-api-service",
    ID:            serviceregistry.GenerateServiceID("my-api-service"),
    Address:       "localhost",
    Port:          8080,
    Tags:          []string{"api", "v1"},
    Meta:          map[string]string{"version": "1.0.0"},
    CheckURL:      "http://localhost:8080/health",
    CheckInterval: "10s",
}

err = serviceregistry.RegisterService(log, serviceInfo)
if err != nil {
    return fmt.Errorf("failed to register service: %w", err)
}
```

### etcd

To use etcd, set `RegistryType=etcd`. etcd supports:

- Service registration with TTL-based health checks
- Service discovery
- Key-based service data storage

Example with etcd:

```go
// Initialize service registry with etcd
config := serviceregistry.Config{
    Address: "localhost:2379",
    Options: map[string]interface{}{
        "username": "my-etcd-user",
        "password": "my-etcd-password",
    },
}

// Initialize registry with etcd
err := serviceregistry.InitServiceRegistry(log, "etcd", config)
if err != nil {
    return fmt.Errorf("failed to initialize service registry: %w", err)
}

// Register service
serviceInfo := serviceregistry.ServiceInfo{
    Name:     "my-api-service",
    ID:       serviceregistry.GenerateServiceID("my-api-service"),
    Address:  "localhost",
    Port:     8080,
    Tags:     []string{"api", "v1"},
    Meta:     map[string]string{"version": "1.0.0"},
    CheckTTL: "30s",
}

err = serviceregistry.RegisterService(log, serviceInfo)
if err != nil {
    return fmt.Errorf("failed to register service: %w", err)
}
```

### Zookeeper

To use Zookeeper, set `RegistryType=zookeeper`. Zookeeper supports:

- Service registration with ephemeral nodes
- Service discovery
- Hierarchical service data storage

Example with Zookeeper:

```go
// Initialize service registry with Zookeeper
config := serviceregistry.Config{
    Address: "localhost:2181",
    Options: map[string]interface{}{
        // Additional servers can be specified
        "servers": []string{"localhost:2181", "localhost:2182", "localhost:2183"},
    },
}

// Initialize registry with Zookeeper
err := serviceregistry.InitServiceRegistry(log, "zookeeper", config)
if err != nil {
    return fmt.Errorf("failed to initialize service registry: %w", err)
}

// Register service
serviceInfo := serviceregistry.ServiceInfo{
    Name:     "my-api-service",
    ID:       serviceregistry.GenerateServiceID("my-api-service"),
    Address:  "localhost",
    Port:     8080,
    Tags:     []string{"api", "v1"},
    Meta:     map[string]string{"version": "1.0.0"},
    CheckTTL: "30s",
}

err = serviceregistry.RegisterService(log, serviceInfo)
if err != nil {
    return fmt.Errorf("failed to register service: %w", err)
}
``` 