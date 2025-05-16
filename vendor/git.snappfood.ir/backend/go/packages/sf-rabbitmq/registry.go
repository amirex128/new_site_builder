package sfrabbitmq

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// =============================================================================
// Registry Types and Global Instance
// =============================================================================

// RegistryOption is a function that configures the Registry
type RegistryOption func(*Registry)

// Registry manages RabbitMQ connections
type Registry struct {
	mu            sync.RWMutex
	connections   map[string]*Connection
	globalOptions []func(*amqp.Config)
	logger        Logger
	outbox        []*OutboxMessage
	retryInterval time.Duration
	maxRetries    int
	declarations  []*Declaration
}

// Connection represents a registered RabbitMQ connection
type Connection struct {
	client     *amqp.Connection
	config     Config
	options    []func(*amqp.Config)
	connecting bool
}

// Global registry instance
var globalRegistry = &Registry{
	connections:   make(map[string]*Connection),
	globalOptions: []func(*amqp.Config){},
	logger:        nil,
	outbox:        make([]*OutboxMessage, 0),
	retryInterval: 5 * time.Second,
	maxRetries:    3,
	declarations:  make([]*Declaration, 0),
}

// RegisterConnection configures RabbitMQ connections with provided options
func RegisterConnection(opts ...RegistryOption) error {
	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	var wg sync.WaitGroup
	wg.Add(len(globalRegistry.connections))

	// Start connections for all registered connections
	for name := range globalRegistry.connections {
		go func(name string) {
			defer wg.Done()
			connectWithRetry(name)
		}(name)
	}

	wg.Wait()

	// Apply all declarations
	return globalRegistry.applyDeclarations()
}

// WithConnectionDetails sets the connection details for a named connection
func WithConnectionDetails(name string, host string, port int, username string, password string, opts ...func(*amqp.Config)) RegistryOption {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*Connection)
		}
		r.connections[name] = &Connection{
			config: Config{
				Host:     host,
				Port:     port,
				Username: username,
				Password: password,
				Vhost:    "/", // Default vhost
			},
			options:    opts,
			connecting: false,
		}
	}
}

// WithOutboxConfig configures the outbox retry behavior
func WithOutboxConfig(retryInterval time.Duration, maxRetries int) RegistryOption {
	return func(r *Registry) {
		r.retryInterval = retryInterval
		r.maxRetries = maxRetries
	}
}

// WithGlobalOptions adds options that will be applied to all connections
func WithGlobalOptions(option func(config *Config)) RegistryOption {
	return func(r *Registry) {
		r.globalOptions = append(r.globalOptions, func(c *amqp.Config) {
			cfg := &Config{}
			option(cfg)
			*c = cfg.ToAMQP()
		})
	}
}

// WithOptions adds connection-specific options
func WithOptions(option func(config *Config)) func(*amqp.Config) {
	return func(c *amqp.Config) {
		cfg := &Config{}
		option(cfg)
		*c = cfg.ToAMQP()
	}
}

// WithLogger sets a custom logger for the registry
func WithLogger(logger Logger) RegistryOption {
	return func(r *Registry) {
		r.logger = logger
	}
}

// WithDeclareExchange declares an exchange with only essential parameters
func WithDeclareExchange(connName string, name string, exchangeType string) RegistryOption {
	return func(r *Registry) {
		r.declarations = append(r.declarations, &Declaration{
			Type:         "exchange",
			ConnName:     connName,
			Exchange:     name,
			ExchangeType: exchangeType,
			Durable:      true,
			AutoDelete:   false,
			Internal:     false,
			NoWait:       false,
		})
	}
}

// WithDeclareQueue declares a queue with only essential parameters
func WithDeclareQueue(connName string, name string) RegistryOption {
	return func(r *Registry) {
		r.declarations = append(r.declarations, &Declaration{
			Type:       "queue",
			ConnName:   connName,
			Queue:      name,
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
		})
	}
}

// WithBind binds a queue to an exchange with only essential parameters
func WithBind(connName string, queueName string, routingKey string, exchangeName string) RegistryOption {
	return func(r *Registry) {
		r.declarations = append(r.declarations, &Declaration{
			Type:       "binding",
			ConnName:   connName,
			Queue:      queueName,
			Exchange:   exchangeName,
			RoutingKey: routingKey,
			NoWait:     false,
		})
	}
}

// WithDeclareExchangeAndQueue declares both exchange and queue and binds them with essential parameters
func WithDeclareExchangeAndQueue(connName string, exchangeName string, exchangeType string, queueName string, routingKey string) RegistryOption {
	return func(r *Registry) {
		// Declare exchange
		WithDeclareExchange(connName, exchangeName, exchangeType)(r)
		// Declare queue
		WithDeclareQueue(connName, queueName)(r)
		// Bind them
		WithBind(connName, queueName, routingKey, exchangeName)(r)
	}
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
				extras := make(map[string]interface{})
				extras["connection_name"] = name
				globalRegistry.logger.InfoWithCategory(Category.Infrastructure.Network, SubCategory.Status.Success, "Successfully connected to RabbitMQ", extras)
			}
			return
		}

		// Connection failed, wait and retry
		if globalRegistry.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = name
			extras[ExtraKey.Error.ErrorMessage] = err.Error()
			extras["retry_interval"] = retryInterval.String()
			globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to connect to RabbitMQ", extras)
		}
		time.Sleep(retryInterval)
	}
}

// initializeConnection attempts to initialize a RabbitMQ connection
func initializeConnection(name string) (*amqp.Connection, error) {
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[name]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// Create AMQP config
	amqpConfig := amqp.Config{
		Vhost: conn.config.Vhost,
	}

	// Apply global options
	for _, gOpt := range globalRegistry.globalOptions {
		gOpt(&amqpConfig)
	}

	// Apply connection-specific options
	for _, opt := range conn.options {
		opt(&amqpConfig)
	}

	// Build connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		conn.config.Username,
		conn.config.Password,
		conn.config.Host,
		conn.config.Port,
		strings.TrimPrefix(conn.config.Vhost, "/"))

	// Create the connection
	client, err := amqp.DialConfig(url, amqpConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ connection: %w", err)
	}

	return client, nil
}

// getConnection returns a RabbitMQ connection for the named connection
func getConnection(name string) (*amqp.Connection, error) {
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

		time.Sleep(100 * time.Millisecond)
	}
}

// Health checks if all RabbitMQ connections are alive
func Health(ctx context.Context) error {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	for name, conn := range globalRegistry.connections {
		if conn.client == nil {
			return fmt.Errorf("connection '%s' is not initialized", name)
		}

		// Check if the connection is still alive
		select {
		case <-conn.client.NotifyClose(make(chan *amqp.Error)):
			return fmt.Errorf("connection '%s' is closed", name)
		default:
			// Connection is still alive
		}
	}

	return nil
}

// CloseConnections closes all RabbitMQ connections
func CloseConnections() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	for name, conn := range globalRegistry.connections {
		if conn.client != nil {
			if err := conn.client.Close(); err != nil {
				if globalRegistry.logger != nil {
					extras := make(map[string]interface{})
					extras["connection_name"] = name
					extras[ExtraKey.Error.ErrorMessage] = err.Error()
					globalRegistry.logger.ErrorWithCategory(Category.Infrastructure.Network, SubCategory.Status.Error, "Failed to close RabbitMQ connection", extras)
				}
			}
		}
		delete(globalRegistry.connections, name)
	}
}

// CloseConnection closes a specific RabbitMQ connection
func CloseConnection(name string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[name]
	if !exists {
		return fmt.Errorf("connection with name '%s' not found", name)
	}

	if conn.client != nil {
		if err := conn.client.Close(); err != nil {
			return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
		}
	}

	delete(globalRegistry.connections, name)
	return nil
}

// applyDeclarations applies all stored declarations
func (r *Registry) applyDeclarations() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, decl := range r.declarations {
		conn, exists := r.connections[decl.ConnName]
		if !exists || conn.client == nil {
			return fmt.Errorf("connection '%s' not available for declaration", decl.ConnName)
		}

		ch, err := conn.client.Channel()
		if err != nil {
			return fmt.Errorf("failed to create channel for declaration: %w", err)
		}
		defer ch.Close()

		switch decl.Type {
		case "exchange":
			err = ch.ExchangeDeclare(
				decl.Exchange,
				decl.ExchangeType,
				decl.Durable,
				decl.AutoDelete,
				decl.Internal,
				decl.NoWait,
				decl.Args,
			)
		case "queue":
			_, err = ch.QueueDeclare(
				decl.Queue,
				decl.Durable,
				decl.AutoDelete,
				decl.Exclusive,
				decl.NoWait,
				decl.Args,
			)
		case "binding":
			err = ch.QueueBind(
				decl.Queue,
				decl.RoutingKey,
				decl.Exchange,
				decl.NoWait,
				decl.Args,
			)
		}

		if err != nil {
			return fmt.Errorf("failed to apply declaration: %w", err)
		}
	}

	return nil
}

// AddToOutbox adds a failed message to the outbox
func (r *Registry) AddToOutbox(msg *OutboxMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	msg.NextRetry = time.Now().Add(r.retryInterval)
	r.outbox = append(r.outbox, msg)
}

// StartOutboxProcessor starts the background worker to process failed messages
func (r *Registry) StartOutboxProcessor() {
	go func() {
		for {
			r.mu.Lock()
			now := time.Now()
			var remaining []*OutboxMessage

			for _, msg := range r.outbox {
				if msg.RetryCount >= r.maxRetries {
					if r.logger != nil {
						extras := make(map[string]interface{})
						extras["exchange"] = msg.Exchange
						extras["routing_key"] = msg.RoutingKey
						extras[ExtraKey.Error.ErrorMessage] = msg.LastError.Error()
						extras["retry_count"] = msg.RetryCount
						r.logger.ErrorWithCategory(Category.Error.Error, SubCategory.Status.Error, "Message permanently failed after max retries", extras)
					}
					continue
				}

				if now.After(msg.NextRetry) {
					// Try to republish
					conn, err := getConnection(msg.ConnectionName)
					if err != nil {
						msg.LastError = err
						msg.RetryCount++
						msg.NextRetry = now.Add(r.retryInterval)
						remaining = append(remaining, msg)
						continue
					}

					ch, err := conn.Channel()
					if err != nil {
						msg.LastError = err
						msg.RetryCount++
						msg.NextRetry = now.Add(r.retryInterval)
						remaining = append(remaining, msg)
						continue
					}

					err = ch.PublishWithContext(
						context.Background(),
						msg.Exchange,
						msg.RoutingKey,
						msg.Mandatory,
						msg.Immediate,
						msg.Message.ToAMQP(),
					)

					ch.Close()

					if err != nil {
						msg.LastError = err
						msg.RetryCount++
						msg.NextRetry = now.Add(r.retryInterval)
						remaining = append(remaining, msg)
						continue
					}

					// Successfully published, remove from outbox
					if r.logger != nil {
						extras := make(map[string]interface{})
						extras["exchange"] = msg.Exchange
						extras["routing_key"] = msg.RoutingKey
						r.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Status.Success, "Successfully republished message from outbox", extras)
					}
				} else {
					remaining = append(remaining, msg)
				}
			}

			r.outbox = remaining
			r.mu.Unlock()

			time.Sleep(r.retryInterval)
		}
	}()
}

// processOutboxMessage attempts to process a message from the outbox
func (r *Registry) processOutboxMessage(msg *OutboxMessage) error {
	// If we've reached the max retries, log an error and remove from outbox
	if msg.RetryCount >= r.maxRetries {
		if r.logger != nil {
			extras := make(map[string]interface{})
			extras["connection_name"] = msg.ConnectionName
			extras["exchange"] = msg.Exchange
			extras["routing_key"] = msg.RoutingKey
			extras["message_id"] = msg.Message.MessageId
			extras["retry_count"] = msg.RetryCount
			extras["max_retries"] = r.maxRetries
			extras[ExtraKey.Error.ErrorMessage] = msg.LastError.Error()
			r.logger.ErrorWithCategory(Category.Error.Error, SubCategory.Status.Error, "Message permanently failed after max retries", extras)
		}
		return msg.LastError
	}

	// ... existing code ...

	// Success - log and return
	if r.logger != nil {
		extras := make(map[string]interface{})
		extras["connection_name"] = msg.ConnectionName
		extras["exchange"] = msg.Exchange
		extras["routing_key"] = msg.RoutingKey
		extras["message_id"] = msg.Message.MessageId
		extras["retry_count"] = msg.RetryCount
		r.logger.InfoWithCategory(Category.API.Messaging, SubCategory.Status.Success, "Successfully republished message from outbox", extras)
	}

	return nil
}
