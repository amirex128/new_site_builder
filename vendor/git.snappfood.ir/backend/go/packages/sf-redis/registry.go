package sfredis

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8/v2"
)

// =============================================================================
// Registry Types and Global Instance
// =============================================================================

// RedisOption is a function that configures a redis.Client instance
type RedisOption func(*Options)

// RegistryOption is a function that configures the Registry
type RegistryOption func(*Registry)

// Registry manages SfRedis connections
type Registry struct {
	mu            sync.RWMutex
	connections   map[string]*Connection
	globalOptions []RedisOption
	logger        Logger
}

// Connection represents a registered SfRedis connection
type Connection struct {
	client     redis.UniversalClient
	addr       string
	password   string
	db         int
	options    []RedisOption
	connecting bool
	config     map[string]interface{}
}

// Global registry instance
var globalRegistry = &Registry{
	connections:   make(map[string]*Connection),
	globalOptions: []RedisOption{},
	logger:        nil,
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
func WithGlobalOptions(option func(*Options)) RegistryOption {
	return func(r *Registry) {
		r.globalOptions = append(r.globalOptions, func(opts *Options) {
			option(opts)
		})
	}
}

// WithConnectionDetails sets the connection details for a named connection
func WithConnectionDetails(name, addr, password string, db int, options ...RedisOption) RegistryOption {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*Connection)
		}
		r.connections[name] = &Connection{
			addr:       addr,
			password:   password,
			db:         db,
			options:    options,
			config:     make(map[string]interface{}),
			connecting: false,
		}
	}
}

// =============================================================================
// Connection Management
// =============================================================================

// RegisterConnection configures SfRedis connections with provided options
func RegisterConnection(opts ...RegistryOption) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
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
			globalRegistry.mu.Unlock()

			if globalRegistry.logger != nil {
				globalRegistry.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Successfully connected to SfRedis", map[string]interface{}{
					"connection_name": name,
				})
			}
			return
		}

		// Connection failed, wait and retry
		if globalRegistry.logger != nil {
			globalRegistry.logger.ErrorWithCategory(Category.Database.Database, Category.Error.Error, "Failed to connect to SfRedis", map[string]interface{}{
				"connection_name": name,
				"error":           err.Error(),
				"retry_interval":  retryInterval.String(),
			})
		}
		time.Sleep(retryInterval)
	}
}

// applyOptions applies the given options to a Redis client
func applyOptions(client redis.UniversalClient, opts *Options) {
	switch c := client.(type) {
	case *redis.Client:
		clientOpts := c.Options()
		clientOpts.PoolSize = opts.PoolSize
		clientOpts.MinIdleConns = opts.MinIdleConns
		clientOpts.MaxRetries = opts.MaxRetries
		clientOpts.DialTimeout = opts.DialTimeout
		clientOpts.ReadTimeout = opts.ReadTimeout
		clientOpts.WriteTimeout = opts.WriteTimeout
		clientOpts.IdleTimeout = opts.IdleTimeout
		clientOpts.MaxConnAge = opts.MaxConnAge
		clientOpts.PoolTimeout = opts.PoolTimeout
	case *redis.ClusterClient:
		clusterOpts := c.Options()
		clusterOpts.PoolSize = opts.PoolSize
		clusterOpts.MinIdleConns = opts.MinIdleConns
		clusterOpts.MaxRetries = opts.MaxRetries
		clusterOpts.DialTimeout = opts.DialTimeout
		clusterOpts.ReadTimeout = opts.ReadTimeout
		clusterOpts.WriteTimeout = opts.WriteTimeout
		clusterOpts.IdleTimeout = opts.IdleTimeout
		clusterOpts.MaxConnAge = opts.MaxConnAge
		clusterOpts.PoolTimeout = opts.PoolTimeout
		clusterOpts.MaxRedirects = opts.ClusterOptions.MaxRedirects
		clusterOpts.RouteRandomly = opts.ClusterOptions.RouteRandomly
		clusterOpts.RouteByLatency = opts.ClusterOptions.RouteByLatency
	}
}

// initializeConnection attempts to initialize a SfRedis connection
// If the address contains commas, it's treated as a cluster configuration
func initializeConnection(name string) (redis.UniversalClient, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	var client redis.UniversalClient

	// Check if address contains commas (cluster configuration)
	if strings.Contains(conn.addr, ",") {
		// Split addresses for cluster nodes
		addresses := strings.Split(conn.addr, ",")
		// Trim spaces from addresses
		for i := range addresses {
			addresses[i] = strings.TrimSpace(addresses[i])
		}

		if globalRegistry.logger != nil {
			globalRegistry.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Initializing Redis Cluster connection", map[string]interface{}{
				"connection_name": name,
				"addresses":       strings.Join(addresses, ","),
			})
		}

		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        addresses,
			Password:     conn.password,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})
	} else {
		if globalRegistry.logger != nil {
			globalRegistry.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Initializing Single Node Redis connection", map[string]interface{}{
				"connection_name": name,
				"address":         conn.addr,
				"database":        conn.db,
			})
		}

		client = redis.NewClient(&redis.Options{
			Addr:     conn.addr,
			Password: conn.password,
			DB:       conn.db,
		})
	}

	// Add APM hook for monitoring
	client.AddHook(apmgoredis.NewHook())

	// Apply global options
	globalOpts := &Options{}
	for _, option := range globalRegistry.globalOptions {
		option(globalOpts)
	}

	applyOptions(client, globalOpts)

	// Apply connection-specific options
	connOpts := &Options{}
	for _, option := range conn.options {
		option(connOpts)
	}
	applyOptions(client, connOpts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// =============================================================================
// Public API
// =============================================================================

// getConnection returns a SfRedis client for the named connection
func getConnection(name string) (redis.UniversalClient, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// If client is already initialized, return it
	if conn.client != nil {
		return conn.client, nil
	}

	// If client is not initialized yet, wait for connection
	for {
		globalRegistry.mu.RLock()
		conn, exists = globalRegistry.connections[name]

		if !exists {
			globalRegistry.mu.RUnlock()
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

// MustClient returns a SfRedis client or panics if it cannot be retrieved
//func MustClient(name string) *redis.Client {
//	client, err := getConnection(name)
//	if err != nil {
//		panic(err)
//	}
//	return client
//}

// Health checks if all SfRedis connections are alive
func Health(ctx context.Context) error {
	globalRegistry.mu.RLock()
	connections := make([]string, 0, len(globalRegistry.connections))
	for name := range globalRegistry.connections {
		connections = append(connections, name)
	}
	globalRegistry.mu.RUnlock()

	if len(connections) == 0 {
		return fmt.Errorf("no SfRedis connections registered")
	}

	var errs []string
	for _, name := range connections {
		client, err := getConnection(name)
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' error: %v", name, err))
			continue
		}

		if err := client.Ping(ctx).Err(); err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' ping error: %v", name, err))
			continue
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("SfRedis health check failed: %s", errs)
	}

	return nil
}

// CloseConnections closes all registered SfRedis connections
func CloseConnections() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	for name, conn := range globalRegistry.connections {
		if conn.client != nil {
			if globalRegistry.logger != nil {
				globalRegistry.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Closing Redis connection", map[string]interface{}{
					"connection_name": name,
				})
			}
			conn.client.Close()
			conn.client = nil
			globalRegistry.connections[name] = conn
		}
	}
}

// CloseConnection closes a specific connection by name
func CloseConnection(name string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[name]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", name)
	}

	if conn.client != nil {
		if globalRegistry.logger != nil {
			globalRegistry.logger.InfoWithCategory(Category.Database.Database, Category.System.Startup, "Closing Redis connection", map[string]interface{}{
				"connection_name": name,
			})
		}
		conn.client.Close()
		conn.client = nil
		globalRegistry.connections[name] = conn
	}

	return nil
}
