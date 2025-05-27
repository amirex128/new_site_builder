package v8

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
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
	retryOptions  *RetryOptions
}

// Connection represents a registered Elasticsearch connection
type Connection struct {
	client       *elasticsearch.Client
	config       elasticsearch.Config
	options      []ClientOption
	connecting   bool
	retryOptions *RetryOptions // Add per-connection retry options
	db           *gorm.DB
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
func WithConnectionDetails(name string, config elasticsearch.Config, options ...interface{}) RegistryOption {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*Connection)
		}
		conn := &Connection{
			config:     config,
			connecting: false,
		}
		for _, opt := range options {
			switch v := opt.(type) {
			case ClientOption:
				conn.options = append(conn.options, v)
			case *RetryOptions:
				conn.retryOptions = v
			}
		}
		r.connections[name] = conn
	}
}

// WithRetryOptions sets custom retry options for the registry
func WithRetryOptions(options *RetryOptions) RegistryOption {
	return func(r *Registry) {
		r.retryOptions = options
	}
}

// =============================================================================
// Connection Management
// =============================================================================

// RegisterConnection configures Elasticsearch connections with provided options
func RegisterConnection(opts ...RegistryOption) error {
	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	// Set default retry options if not set
	if globalRegistry.retryOptions == nil {
		globalRegistry.retryOptions = DefaultRetryOptions()
	}

	var errs []error
	// Start connections for all registered connections
	for name := range globalRegistry.connections {
		err := connectWithRetry(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("connection '%s': %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to connect to one or more databases: %v", errs)
	}
	return nil
}

// connectWithRetry attempts to establish a connection with retry and returns error if all retries fail
func connectWithRetry(name string) error {
	operation := func() error {
		globalRegistry.mu.RLock()
		conn, exists := globalRegistry.connections[name]
		globalRegistry.mu.RUnlock()
		if !exists || conn.db != nil {
			return nil // Already connected or removed
		}
		// Mark that we're connecting
		globalRegistry.mu.Lock()
		if !conn.connecting {
			conn.connecting = true
			globalRegistry.connections[name] = conn
		}
		globalRegistry.mu.Unlock()

		db, err := initializeConnection(name)
		if err == nil && db != nil {
			globalRegistry.mu.Lock()
			conn, stillExists := globalRegistry.connections[name]
			if stillExists {
				conn.db = db
				conn.connecting = false
				globalRegistry.connections[name] = conn
			}
			globalRegistry.mu.Unlock()
			if globalRegistry.logger != nil {
				extras := map[string]interface{}{
					ExtraKey.Database.Table: name,
				}
				globalRegistry.logger.InfoWithCategory(Category.Database.Database, SubCategory.Status.Success, "Successfully connected to database", extras)
			}
			return nil
		}
		return err
	}
	return WithRetry(
		operation,
		globalRegistry.retryOptions,
		globalRegistry.logger,
		"connectWithRetry("+name+")",
	)
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

	// Check client health after creation
	res, err := client.Cluster.Health()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch client health check failed: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch client health check returned error status: %s", res.String())
	}

	return client, nil
}

// =============================================================================
// Public API
// =============================================================================

// GetConnection returns a DB instance for the named connection
func GetConnection(name string) (*gorm.DB, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// If DB is already initialized, return it
	if conn.db != nil {
		return conn.db, nil
	}

	// If DB is not initialized yet, wait for connection
	for {
		globalRegistry.mu.RLock()
		conn, exists = globalRegistry.connections[name]

		if !exists {
			globalRegistry.mu.RUnlock()
			return nil, fmt.Errorf("connection with name '%s' not registered", name)
		}

		if conn.db != nil {
			db := conn.db
			globalRegistry.mu.RUnlock()
			return db, nil
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
	if len(connections) == 0 {
		if globalRegistry.logger != nil {
			globalRegistry.logger.WarnWithCategory(
				Category.System.General,
				SubCategory.Status.Warning,
				"No Elasticsearch connections registered",
				nil,
			)
		}
		globalRegistry.mu.RUnlock()
		return fmt.Errorf("no Elasticsearch connections registered")
	}

	var errs []string
	for _, name := range connections {
		client := MustClient(name)

		// Check cluster health
		res, err := client.Cluster.Health()
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' health check error: %v", name, err))
			if globalRegistry.logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name
				extraMap["error"] = err.Error()

				globalRegistry.logger.ErrorWithCategory(
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
			if globalRegistry.logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name
				extraMap["response"] = res.String()

				globalRegistry.logger.ErrorWithCategory(
					Category.System.Health,
					SubCategory.Status.Error,
					"Elasticsearch connection responded with error status",
					extraMap,
				)
			}
			continue
		}

		if globalRegistry.logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			globalRegistry.logger.InfoWithCategory(
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

	globalRegistry.mu.RUnlock()
	return nil
}

// CloseConnections closes all registered Elasticsearch connections
func CloseConnections() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	for name, conn := range globalRegistry.connections {
		if conn.client != nil {
			// Elasticsearch client doesn't have a Close method
			// But we can set it to nil to indicate it's closed
			conn.client = nil
			globalRegistry.connections[name] = conn
			if globalRegistry.logger != nil {
				// Convert category types to strings for logging
				extraMap := make(map[string]interface{})
				extraMap["connection"] = name

				globalRegistry.logger.InfoWithCategory(
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

	conn, exists := globalRegistry.connections[name]
	if !exists {
		if globalRegistry.logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			globalRegistry.logger.ErrorWithCategory(
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
		if globalRegistry.logger != nil {
			// Convert category types to strings for logging
			extraMap := make(map[string]interface{})
			extraMap["connection"] = name

			globalRegistry.logger.InfoWithCategory(
				Category.System.General,
				SubCategory.Networking.Disconnection,
				"Closed Elasticsearch connection",
				extraMap,
			)
		}
	}

	return nil
}
