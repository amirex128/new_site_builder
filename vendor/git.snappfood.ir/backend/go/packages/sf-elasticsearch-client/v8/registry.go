package v8

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

// =============================================================================
// Registry Types and Global Instance
// =============================================================================

// ClientOption is a function that configures an elasticsearch.Client instance
type ClientOption func(*elasticsearch.Client)

// RegistryOption is a function that configures the Registry
type RegistryOption func(*Registry)

// Registry manages Elasticsearch client connections
type Registry struct {
	mu            sync.RWMutex
	connections   map[string]*Connection
	globalOptions []ClientOption
	logger        Logger
}

// Connection represents a registered Elasticsearch connection
type Connection struct {
	client     *elasticsearch.Client
	config     elasticsearch.Config
	options    []ClientOption
	connecting bool
}

// Global registry instance
var globalRegistry = &Registry{
	connections:   make(map[string]*Connection),
	globalOptions: []ClientOption{},
	logger:        nil, // Will be set during RegisterConnection
}

// =============================================================================
// Registry Configuration Options
// =============================================================================

// WithLogger sets a custom logger for the registry
func WithLogger(logger Logger) RegistryOption {
	return func(r *Registry) {
		r.logger = logger
	}
}

// WithGlobalOptions adds options that will be applied to all new connections
func WithGlobalOptions(option func(*elasticsearch.Client)) RegistryOption {
	return func(r *Registry) {
		r.globalOptions = append(r.globalOptions, option)
	}
}

// WithConnectionDetails sets the connection details for a named connection
func WithConnectionDetails(name string, config elasticsearch.Config, options ...ClientOption) RegistryOption {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*Connection)
		}
		r.connections[name] = &Connection{
			config:     config,
			options:    options,
			connecting: false,
		}
	}
}

// =============================================================================
// Connection Management
// =============================================================================

// RegisterConnection configures Elasticsearch connections with provided options
func RegisterConnection(opts ...RegistryOption) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	// Ensure logger is set
	if globalRegistry.logger == nil {
		return fmt.Errorf("logger must be provided using WithLogger option")
	}

	// Start connections for all registered connections
	for name := range globalRegistry.connections {
		go connectWithRetry(name)
	}

	return nil
}

// connectWithRetry attempts to establish a connection with infinite retry
func connectWithRetry(name string) {
	retryInterval := 5 * time.Second

	for {
		// Check if we're still supposed to be connecting
		globalRegistry.mu.RLock()
		conn, exists := globalRegistry.connections[name]
		logger := globalRegistry.logger
		if !exists || conn.client != nil {
			globalRegistry.mu.RUnlock()
			return
		}

		// Mark that we're connecting
		if !conn.connecting {
			globalRegistry.mu.RUnlock()
			globalRegistry.mu.Lock()
			conn.connecting = true
			globalRegistry.mu.Unlock()
			globalRegistry.mu.RLock()
		}
		globalRegistry.mu.RUnlock()

		// Try to connect
		client, err := initializeConnection(name)

		if err == nil && client != nil {
			// Connection successful
			globalRegistry.mu.Lock()
			conn, stillExists := globalRegistry.connections[name]
			if stillExists {
				conn.client = client
				conn.connecting = false
				globalRegistry.connections[name] = conn
			}
			logger := globalRegistry.logger
			globalRegistry.mu.Unlock()

			if logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name

				logger.InfoWithCategory(
					Category.System.General,
					SubCategory.Networking.Connection,
					"Successfully connected to Elasticsearch",
					extraMap,
				)
			}
			return
		}

		// Connection failed, wait and retry
		if logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name
			extraMap["error"] = err.Error()
			extraMap["retryInterval"] = retryInterval.String()

			logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Networking.Connection,
				"Failed to connect to Elasticsearch",
				extraMap,
			)
		}
		time.Sleep(retryInterval)
	}
}

// initializeConnection attempts to initialize an Elasticsearch connection
func initializeConnection(name string) (*elasticsearch.Client, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// Create the client
	client, err := elasticsearch.NewClient(conn.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Apply global options
	for _, option := range globalRegistry.globalOptions {
		option(client)
	}

	// Apply connection-specific options
	for _, option := range conn.options {
		option(client)
	}

	return client, nil
}

// =============================================================================
// Public API
// =============================================================================

// GetConnection returns an Elasticsearch client for the named connection
func GetConnection(name string) (*elasticsearch.Client, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	logger := globalRegistry.logger
	globalRegistry.mu.RUnlock()

	if !exists {
		if logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Networking.Connection,
				"Connection not registered",
				extraMap,
			)
		}
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// If client is already initialized, return it
	if conn.client != nil {
		return conn.client, nil
	}

	// If client is not initialized yet, wait for connection
	if logger != nil {
		// Convert category types to strings for logging
		extraMap := make(map[string]interface{})
		extraMap["connection"] = name

		logger.InfoWithCategory(
			Category.System.General,
			SubCategory.Networking.Connection,
			"Waiting for connection to initialize",
			extraMap,
		)
	}

	for {
		globalRegistry.mu.RLock()
		conn, exists = globalRegistry.connections[name]

		if !exists {
			globalRegistry.mu.RUnlock()
			if logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name

				logger.ErrorWithCategory(
					Category.System.General,
					SubCategory.Networking.Connection,
					"Connection disappeared during initialization",
					extraMap,
				)
			}
			return nil, fmt.Errorf("connection with name '%s' not registered", name)
		}

		if conn.client != nil {
			client := conn.client
			globalRegistry.mu.RUnlock()
			return client, nil
		}

		// If not connecting yet, start connecting
		if !conn.connecting {
			globalRegistry.mu.RUnlock()
			go connectWithRetry(name)
		} else {
			globalRegistry.mu.RUnlock()
		}

		// Wait a bit before checking again
		time.Sleep(100 * time.Millisecond)
	}
}

// MustClient returns an Elasticsearch client or panics if it cannot be retrieved
func MustClient(name string) *elasticsearch.Client {
	client, err := GetConnection(name)
	if err != nil {
		panic(err)
	}
	return client
}

// Health checks if all Elasticsearch connections are alive
func Health(ctx context.Context) error {
	globalRegistry.mu.RLock()
	connections := make([]string, 0, len(globalRegistry.connections))
	for name := range globalRegistry.connections {
		connections = append(connections, name)
	}
	logger := globalRegistry.logger
	globalRegistry.mu.RUnlock()

	if len(connections) == 0 {
		if logger != nil {
			logger.WarnWithCategory(
				Category.System.General,
				SubCategory.Status.Warning,
				"No Elasticsearch connections registered",
				nil,
			)
		}
		return fmt.Errorf("no Elasticsearch connections registered")
	}

	var errs []string
	for _, name := range connections {
		client := MustClient(name)

		// Check cluster health
		res, err := client.Cluster.Health()
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' health check error: %v", name, err))
			if logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name
				extraMap["error"] = err.Error()

				logger.ErrorWithCategory(
					Category.System.Health,
					SubCategory.Status.Error,
					"Elasticsearch connection health check failed",
					extraMap,
				)
			}
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			errs = append(errs, fmt.Sprintf("connection '%s' health check error: %s", name, res.String()))
			if logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name
				extraMap["response"] = res.String()

				logger.ErrorWithCategory(
					Category.System.Health,
					SubCategory.Status.Error,
					"Elasticsearch connection responded with error status",
					extraMap,
				)
			}
			continue
		}

		if logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			logger.InfoWithCategory(
				Category.System.Health,
				SubCategory.Status.Success,
				"Elasticsearch connection health check passed",
				extraMap,
			)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Elasticsearch health check failed: %s", strings.Join(errs, "; "))
	}

	return nil
}

// CloseConnections closes all registered Elasticsearch connections
func CloseConnections() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	logger := globalRegistry.logger
	for name, conn := range globalRegistry.connections {
		if conn.client != nil {
			// Elasticsearch client doesn't have a Close method
			// But we can set it to nil to indicate it's closed
			conn.client = nil
			globalRegistry.connections[name] = conn
			if logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name

				logger.InfoWithCategory(
					Category.System.General,
					SubCategory.Networking.Disconnection,
					"Closed Elasticsearch connection",
					extraMap,
				)
			}
		}
	}
}

// CloseConnection closes a specific connection by name
func CloseConnection(name string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	logger := globalRegistry.logger
	conn, exists := globalRegistry.connections[name]
	if !exists {
		if logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			logger.ErrorWithCategory(
				Category.System.General,
				SubCategory.Networking.Disconnection,
				"Failed to close connection: not registered",
				extraMap,
			)
		}
		return fmt.Errorf("connection with name '%s' not registered", name)
	}

	if conn.client != nil {
		// Elasticsearch client doesn't have a Close method
		// But we can set it to nil to indicate it's closed
		conn.client = nil
		globalRegistry.connections[name] = conn
		if logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			logger.InfoWithCategory(
				Category.System.General,
				SubCategory.Networking.Disconnection,
				"Closed Elasticsearch connection",
				extraMap,
			)
		}
	}

	return nil
}
