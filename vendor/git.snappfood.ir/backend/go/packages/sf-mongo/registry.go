//go:build go1.18
// +build go1.18

package sfmongo

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// =============================================================================
// Registry Types and Global Instance
// =============================================================================

// MongoClientOption is a function that configures a mongo.Client instance
// (after creation, before use)
type MongoClientOption func(*mongo.Client)

// MongoRegistryOption is a function that configures the MongoRegistry
// (for global options, logger, etc)
type MongoRegistryOption func(*MongoRegistry)

// MongoRegistry manages MongoDB client connections
type MongoRegistry struct {
	mu            sync.RWMutex
	connections   map[string]*MongoConnection
	globalOptions []MongoClientOption
	logger        Logger
	retryOptions  *RetryOptions
}

// MongoConnection represents a registered MongoDB connection
type MongoConnection struct {
	client       *mongo.Client
	config       *options.ClientOptions
	options      []MongoClientOption
	connecting   bool
	retryOptions *RetryOptions // per-connection retry options
}

// Global MongoDB registry instance
var globalMongoRegistry = &MongoRegistry{
	connections:   make(map[string]*MongoConnection),
	globalOptions: []MongoClientOption{},
	logger:        nil, // Set during RegisterConnection
}

// =============================================================================
// Registry Configuration Options
// =============================================================================

// WithLogger sets a custom logger for the MongoDB registry
func WithLogger(logger Logger) MongoRegistryOption {
	return func(r *MongoRegistry) {
		r.logger = logger
	}
}

// WithGlobalOptions adds options that will be applied to all new MongoDB connections
func WithGlobalOptions(option MongoClientOption) MongoRegistryOption {
	return func(r *MongoRegistry) {
		r.globalOptions = append(r.globalOptions, option)
	}
}

// WithConnectionDetails sets the connection details for a named MongoDB connection
func WithConnectionDetails(name string, config *options.ClientOptions, options ...interface{}) MongoRegistryOption {
	return func(r *MongoRegistry) {
		if r.connections == nil {
			r.connections = make(map[string]*MongoConnection)
		}
		conn := &MongoConnection{
			config:     config,
			connecting: false,
		}
		for _, opt := range options {
			switch v := opt.(type) {
			case MongoClientOption:
				conn.options = append(conn.options, v)
			case *RetryOptions:
				conn.retryOptions = v
			}
		}
		r.connections[name] = conn
	}
}

// WithRetryOptions sets custom retry options for the MongoDB registry
func WithRetryOptions(options *RetryOptions) MongoRegistryOption {
	return func(r *MongoRegistry) {
		r.retryOptions = options
	}
}

// =============================================================================
// Connection Management
// =============================================================================

// RegisterConnection configures MongoDB connections with provided options
func RegisterConnection(opts ...MongoRegistryOption) error {
	globalMongoRegistry.mu.Lock()
	defer globalMongoRegistry.mu.Unlock()

	// Apply options
	for _, opt := range opts {
		opt(globalMongoRegistry)
	}

	// Ensure logger is set
	if globalMongoRegistry.logger == nil {
		return fmt.Errorf("logger must be provided using WithLogger option")
	}

	// Start connections for all registered connections
	for name := range globalMongoRegistry.connections {
		go connectWithRetry(name)
	}

	return nil
}

// connectWithRetry attempts to establish a connection with retry and backoff
func connectWithRetry(name string) {
	operation := func() error {
		globalMongoRegistry.mu.RLock()
		conn, exists := globalMongoRegistry.connections[name]
		if !exists || conn.client != nil {
			globalMongoRegistry.mu.RUnlock()
			return nil
		}
		if !conn.connecting {
			globalMongoRegistry.mu.RUnlock()
			globalMongoRegistry.mu.Lock()
			conn.connecting = true
			globalMongoRegistry.mu.Unlock()
			globalMongoRegistry.mu.RLock()
		}
		globalMongoRegistry.mu.RUnlock()

		client, err := initializeMongoConnection(name)
		if err == nil && client != nil {
			globalMongoRegistry.mu.Lock()
			conn, stillExists := globalMongoRegistry.connections[name]
			if stillExists {
				conn.client = client
				conn.connecting = false
				globalMongoRegistry.connections[name] = conn
			}
			globalMongoRegistry.mu.Unlock()

			if globalMongoRegistry.logger != nil {
				extraMap := map[string]interface{}{"connection": name}
				globalMongoRegistry.logger.InfoWithCategory(
					Category.Database.MongoDB,
					SubCategory.Networking.Connection,
					"Successfully connected to MongoDB",
					extraMap,
				)
			}
			return nil
		}
		return err
	}

	globalMongoRegistry.mu.RLock()
	conn, exists := globalMongoRegistry.connections[name]
	var retryOpts *RetryOptions
	if exists && conn.retryOptions != nil {
		retryOpts = conn.retryOptions
	} else {
		retryOpts = globalMongoRegistry.retryOptions
	}
	globalMongoRegistry.mu.RUnlock()
	WithRetry(operation, retryOpts, globalMongoRegistry.logger, "connectWithRetry: "+name)
}

// initializeMongoConnection attempts to initialize a MongoDB connection
func initializeMongoConnection(name string) (*mongo.Client, error) {
	globalMongoRegistry.mu.RLock()
	conn, exists := globalMongoRegistry.connections[name]
	globalMongoRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	client, err := mongo.Connect(context.Background(), conn.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	for _, option := range globalMongoRegistry.globalOptions {
		option(client)
	}
	for _, option := range conn.options {
		option(client)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongodb client ping failed: %w", err)
	}

	return client, nil
}

// =============================================================================
// Public API
// =============================================================================

// GetConnection returns a MongoDB client for the named connection
func GetConnection(name string) (*mongo.Client, error) {
	globalMongoRegistry.mu.RLock()
	conn, exists := globalMongoRegistry.connections[name]
	if !exists {
		if globalMongoRegistry.logger != nil {
			extraMap := map[string]interface{}{"connection": name}
			globalMongoRegistry.logger.ErrorWithCategory(
				Category.Database.MongoDB,
				SubCategory.Networking.Connection,
				"Connection not registered",
				extraMap,
			)
		}
		globalMongoRegistry.mu.RUnlock()
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}
	if conn.client != nil {
		globalMongoRegistry.mu.RUnlock()
		return conn.client, nil
	}
	if globalMongoRegistry.logger != nil {
		extraMap := map[string]interface{}{"connection": name}
		globalMongoRegistry.logger.InfoWithCategory(
			Category.Database.MongoDB,
			SubCategory.Networking.Connection,
			"Waiting for connection to initialize",
			extraMap,
		)
	}
	for {
		globalMongoRegistry.mu.RLock()
		conn, exists = globalMongoRegistry.connections[name]
		if !exists {
			globalMongoRegistry.mu.RUnlock()
			if globalMongoRegistry.logger != nil {
				extraMap := map[string]interface{}{"connection": name}
				globalMongoRegistry.logger.ErrorWithCategory(
					Category.Database.MongoDB,
					SubCategory.Networking.Connection,
					"Connection disappeared during initialization",
					extraMap,
				)
			}
			return nil, fmt.Errorf("connection with name '%s' not registered", name)
		}
		if conn.client != nil {
			client := conn.client
			globalMongoRegistry.mu.RUnlock()
			return client, nil
		}
		if !conn.connecting {
			globalMongoRegistry.mu.RUnlock()
			go connectWithRetry(name)
		} else {
			globalMongoRegistry.mu.RUnlock()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// MustMongoClient returns a MongoDB client or panics if it cannot be retrieved
func MustMongoClient(name string) *mongo.Client {
	client, err := GetConnection(name)
	if err != nil {
		panic(err)
	}
	return client
}

// Health checks if all MongoDB connections are alive
func Health(ctx context.Context) error {
	globalMongoRegistry.mu.RLock()
	connections := make([]string, 0, len(globalMongoRegistry.connections))
	for name := range globalMongoRegistry.connections {
		connections = append(connections, name)
	}
	if len(connections) == 0 {
		if globalMongoRegistry.logger != nil {
			globalMongoRegistry.logger.WarnWithCategory(
				Category.Database.MongoDB,
				SubCategory.Status.Warning,
				"No MongoDB connections registered",
				nil,
			)
		}
		globalMongoRegistry.mu.RUnlock()
		return fmt.Errorf("no MongoDB connections registered")
	}

	var errs []string
	for _, name := range connections {
		client := MustMongoClient(name)
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := client.Ping(pingCtx, nil)
		cancel()
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' health check error: %v", name, err))
			if globalMongoRegistry.logger != nil {
				extraMap := map[string]interface{}{"connection": name, "error": err.Error()}
				globalMongoRegistry.logger.ErrorWithCategory(
					Category.Database.MongoDB,
					SubCategory.Status.Error,
					"MongoDB connection health check failed",
					extraMap,
				)
			}
			continue
		}
		if globalMongoRegistry.logger != nil {
			extraMap := map[string]interface{}{"connection": name}
			globalMongoRegistry.logger.InfoWithCategory(
				Category.Database.MongoDB,
				SubCategory.Status.Success,
				"MongoDB connection health check passed",
				extraMap,
			)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("MongoDB health check failed: %s", strings.Join(errs, "; "))
	}
	globalMongoRegistry.mu.RUnlock()
	return nil
}

// CloseConnections closes all registered MongoDB connections
func CloseConnections() {
	globalMongoRegistry.mu.Lock()
	defer globalMongoRegistry.mu.Unlock()
	for name, conn := range globalMongoRegistry.connections {
		if conn.client != nil {
			_ = conn.client.Disconnect(context.Background())
			conn.client = nil
			globalMongoRegistry.connections[name] = conn
			if globalMongoRegistry.logger != nil {
				extraMap := map[string]interface{}{"connection": name}
				globalMongoRegistry.logger.InfoWithCategory(
					Category.Database.MongoDB,
					SubCategory.Networking.Disconnection,
					"Closed MongoDB connection",
					extraMap,
				)
			}
		}
	}
}

// CloseConnection closes a specific MongoDB connection by name
func CloseConnection(name string) error {
	globalMongoRegistry.mu.Lock()
	defer globalMongoRegistry.mu.Unlock()
	conn, exists := globalMongoRegistry.connections[name]
	if !exists {
		if globalMongoRegistry.logger != nil {
			extraMap := map[string]interface{}{"connection": name}
			globalMongoRegistry.logger.ErrorWithCategory(
				Category.Database.MongoDB,
				SubCategory.Networking.Disconnection,
				"Failed to close connection: not registered",
				extraMap,
			)
		}
		return fmt.Errorf("connection with name '%s' not registered", name)
	}
	if conn.client != nil {
		_ = conn.client.Disconnect(context.Background())
		conn.client = nil
		globalMongoRegistry.connections[name] = conn
		if globalMongoRegistry.logger != nil {
			extraMap := map[string]interface{}{"connection": name}
			globalMongoRegistry.logger.InfoWithCategory(
				Category.Database.MongoDB,
				SubCategory.Networking.Disconnection,
				"Closed MongoDB connection",
				extraMap,
			)
		}
	}
	return nil
}
