# Configuration Manager Package

The `sfconfigmanager` package provides a flexible configuration system for Go applications, supporting multiple configuration sources with a unified interface.

## Features

- **Multiple Backend Support**: Load configuration from various sources including:
  - Vault
  - Environment Variables
  - Configuration Files (JSON, YAML)
  - Consul
  - etcd
  - ZooKeeper

- **Fallback Support**: Define priority order for configuration sources
- **Type Conversion**: Built-in helpers for retrieving config values
- **Struct Binding**: Load configuration directly into Go structs
- **Retry Mechanism**: Built-in retry with exponential backoff for resilient configuration loading

## Installation

```bash
go get github.com/yourdomain/sf-config-manager
```

## Usage

### Basic Usage

```go
package main

import (
	"log"
	"time"
	
	"github.com/yourdomain/sfconfigmanager"
)

// Define your configuration structure
type Config struct {
	Server struct {
		Host    string `env:"SERVER_HOST"`
		Port    int    `env:"SERVER_PORT"`
		Timeout int    `env:"SERVER_TIMEOUT"`
	}
	Database struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
	}
	Cache struct {
		Host     string `env:"CACHE_HOST"`
		Port     int    `env:"CACHE_PORT"`
		Password string `env:"CACHE_PASSWORD"`
		DB       int    `env:"CACHE_DB"`
	}
}

func main() {
	// Initialize a configuration struct
	config := &Config{}
	
	// Register connection with fallback sources
	config, err := sfconfigmanager.RegisterConnection(
		sfconfigmanager.WithConfig(config),
		sfconfigmanager.WithVaultOptions("http://vault:8200", "vault-token", &sfconfigmanager.VaultOptions{
			SecretPath: "app/config",
			Timeout: 5 * time.Second,
		}),
		sfconfigmanager.WithFileOptions("config.json", &sfconfigmanager.FileOptions{
			Type: "json",
		}),
	)
	
	if err != nil {
		log.Fatalf("Failed to register config manager connections: %v", err)
	}
	
	// Now you can use the loaded configuration
	log.Printf("Server host: %s, port: %d", config.Server.Host, config.Server.Port)
}
```

### Using Multiple Config Sources with Fallback

By specifying multiple connection options, you can define a fallback strategy. The library will try each source in the order they are specified until one succeeds:

```go
config, err := sfconfigmanager.RegisterConnection(
	sfconfigmanager.WithConfig(config),
	// Try to load from Vault first
	sfconfigmanager.WithVaultOptions("http://vault:8200", "vault-token", &sfconfigmanager.VaultOptions{
		SecretPath: "app/config",
	}),
	// If Vault fails, try to load from Consul
	sfconfigmanager.WithConsulOptions("localhost:8500", &sfconfigmanager.ConsulOptions{
		Username: "user",
		Password: "pass",
		Prefix: "app/config",
	}),
	// If Consul fails, try to load from environment variables
	sfconfigmanager.WithEnvOptions(&sfconfigmanager.EnvOptions{
		Prefix: "APP_",
	}),
	// As a last resort, try to load from a local file
	sfconfigmanager.WithFileOptions("config.json", nil),
)
```

### Configuring Retry Behavior

You can configure the retry behavior for connection attempts to each configuration source:

```go
config, err := sfconfigmanager.RegisterConnection(
	sfconfigmanager.WithConfig(config),
	// Configure retry options
	sfconfigmanager.WithRetryOptions(&sfconfigmanager.RetryOptions{
		MaxRetries:     5,            // Maximum number of retry attempts
		InitialBackoff: time.Second,  // Initial waiting time between retries
		MaxBackoff:     15 * time.Second, // Maximum waiting time between retries
		BackoffFactor:  1.5,         // Exponential backoff multiplier
	}),
	// Now define your configuration sources
	sfconfigmanager.WithVaultOptions("http://vault:8200", "vault-token", nil),
	sfconfigmanager.WithFileOptions("config.json", nil),
)
```

If no retry options are provided, the system uses the following defaults:
- 3 maximum retries
- 500ms initial backoff
- 10s maximum backoff
- 2.0 backoff factor (doubles the wait time after each failure)

Each configuration provider will be retried according to the specified options before moving on to the next provider in the fallback chain.

### Available Configuration Options

#### Vault Options

```go
// VaultOptions contains Vault-specific configuration options
type VaultOptions struct {
	Address     string        // Vault server address
	Token       string        // Vault authentication token
	SecretPath  string        // Path to the secret
	SecretMount string        // Secret engine mount point
	Timeout     time.Duration // Request timeout
	AuthMethod  string        // Authentication method (token, kubernetes, etc.)
	Role        string        // Role for authentication
}

// Using Vault with options
sfconfigmanager.WithVaultOptions("http://vault:8200", "vault-token", &sfconfigmanager.VaultOptions{
	SecretPath: "secret/data/myapp",
	SecretMount: "secret",
	Timeout: 10 * time.Second,
})
```

#### File Options

```go
// FileOptions contains File-specific configuration options
type FileOptions struct {
	Path string // Path to the configuration file
	Type string // File type: "json", "yaml", or "yml"
}

// Using File with options
sfconfigmanager.WithFileOptions("config.json", &sfconfigmanager.FileOptions{
	Type: "json", // or "yaml"
})
```

#### Environment Options

```go
// EnvOptions contains Environment-specific configuration options
type EnvOptions struct {
	Prefix string // Prefix for environment variables
}

// Using environment variables with options
sfconfigmanager.WithEnvOptions(&sfconfigmanager.EnvOptions{
	Prefix: "MYAPP_", // Will match MYAPP_* environment variables
})
```

#### Consul Options

```go
// ConsulOptions contains Consul-specific configuration options
type ConsulOptions struct {
	Address    string        // Consul server address
	Datacenter string        // Datacenter name
	Token      string        // ACL token
	Timeout    time.Duration // Request timeout
	Prefix     string        // Key prefix
	Username   string        // Basic auth username
	Password   string        // Basic auth password
}

// Using Consul with options
sfconfigmanager.WithConsulOptions("localhost:8500", &sfconfigmanager.ConsulOptions{
	Prefix: "myapp/config/",
	Token: "consul-token",
	Timeout: 5 * time.Second,
})
```

#### etcd Options

```go
// EtcdOptions contains etcd-specific configuration options
type EtcdOptions struct {
	Endpoints []string      // List of etcd endpoints
	Username  string        // Authentication username
	Password  string        // Authentication password
	Timeout   time.Duration // Request timeout
	BasePath  string        // Key prefix
}

// Using etcd with options
sfconfigmanager.WithEtcdOptions([]string{"localhost:2379"}, &sfconfigmanager.EtcdOptions{
	BasePath: "/myapp/config",
	Timeout: 5 * time.Second,
})
```

#### Zookeeper Options

```go
// ZookeeperOptions contains Zookeeper-specific configuration options
type ZookeeperOptions struct {
	Servers  []string      // List of Zookeeper servers
	Timeout  time.Duration // Session timeout
	BasePath string        // Base path for configuration
	Username string        // Authentication username
	Password string        // Authentication password
}

// Using Zookeeper with options
sfconfigmanager.WithZookeeperOptions([]string{"localhost:2181"}, &sfconfigmanager.ZookeeperOptions{
	BasePath: "/myapp/config",
	Timeout: 10 * time.Second,
})
```

### Adding a Custom Logger

You can customize logger behavior by implementing the Logger interface:

```go
import (
	"os"
	
	"github.com/yourdomain/sfconfigmanager"
	"github.com/yourimplementedlogger" // Your logger implementation
)

// Create a logger that implements sfconfigmanager.Logger interface
logger := yourimplementedlogger.New(os.Stdout)

config, err := sfconfigmanager.RegisterConnection(
	sfconfigmanager.WithConfig(config),
	sfconfigmanager.WithLogger(logger),
	sfconfigmanager.WithVaultOptions("vault-addr", "vault-token", nil),
)
```

### Backward Compatibility

For backward compatibility, you can still use the `WithConnectionDetails` function:

```go
config, err := sfconfigmanager.RegisterConnection(
	sfconfigmanager.WithConfig(config),
	sfconfigmanager.WithConnectionDetails(sfconfigmanager.VaultManager, "http://vault:8200", "", "vault-token", "app/config"),
	sfconfigmanager.WithConnectionDetails(sfconfigmanager.FileManager, "config.json", "", ""),
)
```

## Struct Tag Format

Use the `env` tag on your struct fields to map configuration values:

```go
type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST"` // Will match SERVER_HOST in environment vars or config
		Port int    `env:"SERVER_PORT"` // Will match SERVER_PORT in environment vars or config
	}
}
``` 