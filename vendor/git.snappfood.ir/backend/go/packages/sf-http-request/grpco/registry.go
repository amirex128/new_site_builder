package grpco

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	sfhttprequest "git.snappfood.ir/backend/go/packages/sf-http-request"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Registry is a global registry for gRPC connections
type Registry struct {
	mu                   sync.RWMutex
	connections          map[string]*registeredConnection
	defaultTransport     grpc.DialOption
	globalHeaders        metadata.MD
	logger               sfhttprequest.Logger
	serviceClients       map[string]interface{}
	autoConnectTimeoutMs int // Timeout for auto-connecting in milliseconds, 0 means disabled
}

// registeredConnection stores connection details and its services
type registeredConnection struct {
	conn              *grpc.ClientConn
	target            string
	dialOptions       []grpc.DialOption
	connectionHeaders metadata.MD
	services          map[string]ServiceDefinition // key -> service definition
}

// Option is a function to customize the Registry
type Option func(*Registry)

// ServiceDefinition defines a gRPC service with its client constructor and methods
type ServiceDefinition struct {
	ClientConstructor interface{}
	Methods           map[string]string // methodKey -> full method path
}

// WithDefaultTransport sets the default transport for all new connections
func WithDefaultTransport(opt grpc.DialOption) Option {
	return func(r *Registry) {
		r.defaultTransport = opt
	}
}

// WithGlobalHeaders sets the global headers for all connections
func WithGlobalHeaders(headers metadata.MD) Option {
	return func(r *Registry) {
		r.globalHeaders = headers
	}
}

// WithServices registers services for a connection
func WithServices(connectionName string, services map[string]ServiceDefinition) Option {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*registeredConnection)
		}

		// Get or create the connection
		conn, exists := r.connections[connectionName]
		if !exists {
			return // Connection doesn't exist, can't register services
		}

		// Initialize services map if needed
		if conn.services == nil {
			conn.services = make(map[string]ServiceDefinition)
		}

		// Add the services to the connection
		for serviceKey, svcDef := range services {
			conn.services[serviceKey] = svcDef

			if r.logger != nil {
				r.logger.InfoWithCategory(
					sfhttprequest.Category.API.GRPC,
					sfhttprequest.SubCategory.Operation.Registration,
					"Service registered",
					map[string]interface{}{
						sfhttprequest.ExtraKey.Service.ServiceName: serviceKey,
					})
			}
		}
	}
}

// WithConnectionDetails sets the connection details for a named connection
// and optionally registers services
func WithConnectionDetails(name, target string, opts ...interface{}) Option {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*registeredConnection)
		}

		// Convert ClientOption to grpc.DialOption
		dialOpts := make([]grpc.DialOption, 0)

		// Create services map to collect all services
		services := make(map[string]ServiceDefinition)

		// Process all options
		for _, opt := range opts {
			switch typedOpt := opt.(type) {
			case ClientOption:
				// It's a dial option
				var dialOpt grpc.DialOption
				typedOpt(&dialOpt)
				dialOpts = append(dialOpts, dialOpt)
			case map[string]ServiceDefinition:
				// It's a services map
				for key, svc := range typedOpt {
					services[key] = svc
				}
			}
		}

		r.connections[name] = &registeredConnection{
			target:            target,
			dialOptions:       dialOpts,
			connectionHeaders: metadata.MD{},
			services:          make(map[string]ServiceDefinition),
		}

		// Add services if any were provided
		if len(services) > 0 {
			conn := r.connections[name]

			// Add the services to the connection
			for serviceKey, svcDef := range services {
				conn.services[serviceKey] = svcDef

				if r.logger != nil {
					r.logger.InfoWithCategory(
						sfhttprequest.Category.API.GRPC,
						sfhttprequest.SubCategory.Operation.Registration,
						"Service registered",
						map[string]interface{}{
							sfhttprequest.ExtraKey.Service.ServiceName: serviceKey,
						})
				}
			}
		}
	}
}

// WithLogger sets the logger for the registry
func WithLogger(logger sfhttprequest.Logger) Option {
	return func(r *Registry) {
		r.logger = logger
	}
}

// WithAutoConnectTimeout sets a timeout for automatic connection waiting
func WithAutoConnectTimeout(timeoutMs int) Option {
	return func(r *Registry) {
		r.autoConnectTimeoutMs = timeoutMs
	}
}

// Global registry instance
var globalRegistry = &Registry{
	connections:          make(map[string]*registeredConnection),
	defaultTransport:     grpc.WithTransportCredentials(insecure.NewCredentials()),
	globalHeaders:        metadata.MD{},
	serviceClients:       make(map[string]interface{}),
	autoConnectTimeoutMs: 0, // Disabled by default
}

// RegisterConnection configures the gRPC connections with the provided options
func RegisterConnection(opts ...Option) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	if globalRegistry.logger != nil {
		globalRegistry.logger.InfoWithCategory(
			sfhttprequest.Category.API.GRPC,
			sfhttprequest.SubCategory.Operation.Registration,
			"gRPC registry configured",
			map[string]interface{}{
				sfhttprequest.ExtraKey.Service.ServiceName: "gRPC registry",
			})
	}

	// Start connections for all registered connections
	for name := range globalRegistry.connections {
		go connectWithRetry(context.Background(), name)
	}

	return nil
}

// connectWithRetry attempts to establish a connection with infinite retry
func connectWithRetry(ctx context.Context, name string) {
	retryInterval := 5 * time.Second

	for {
		// Check if we're still supposed to be connecting
		globalRegistry.mu.RLock()
		regConn, exists := globalRegistry.connections[name]
		logger := globalRegistry.logger
		if !exists || regConn.conn != nil {
			globalRegistry.mu.RUnlock()
			return
		}
		globalRegistry.mu.RUnlock()

		// Try to connect
		conn, err := initializeConnection(name)

		if err == nil && conn != nil {
			// Connection successful
			globalRegistry.mu.Lock()
			regConn, stillExists := globalRegistry.connections[name]
			if stillExists {
				regConn.conn = conn
				globalRegistry.connections[name] = regConn

				if logger != nil {
					logger.InfoWithCategory(
						sfhttprequest.Category.API.GRPC,
						sfhttprequest.SubCategory.Networking.Connection,
						"gRPC connection established",
						map[string]interface{}{
							sfhttprequest.ExtraKey.Service.ServiceName: name,
							sfhttprequest.ExtraKey.HTTP.URL:            regConn.target,
						})
				}
			}
			globalRegistry.mu.Unlock()
			return
		}

		// Connection failed, log the error
		if logger != nil {
			logger.ErrorWithCategory(
				sfhttprequest.Category.API.GRPC,
				sfhttprequest.SubCategory.Status.Error,
				"Failed to establish gRPC connection",
				map[string]interface{}{
					sfhttprequest.ExtraKey.Service.ServiceName: name,
					sfhttprequest.ExtraKey.Error.ErrorMessage:  err.Error(),
				})
		}

		// Connection failed, wait and retry
		time.Sleep(retryInterval)
	}
}

// initializeConnection creates a new gRPC client instance
func initializeConnection(name string) (*grpc.ClientConn, error) {
	globalRegistry.mu.RLock()
	regConn, exists := globalRegistry.connections[name]
	logger := globalRegistry.logger
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", name)
	}

	// Ensure we have at least the default transport option
	if len(regConn.dialOptions) == 0 {
		regConn.dialOptions = append(regConn.dialOptions, globalRegistry.defaultTransport)
	}

	// Add APM interceptors
	regConn.dialOptions = append(regConn.dialOptions,
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
	)

	if logger != nil {
		logger.DebugWithCategory(
			sfhttprequest.Category.API.GRPC,
			sfhttprequest.SubCategory.Networking.Connection,
			"Initializing gRPC connection",
			map[string]interface{}{
				sfhttprequest.ExtraKey.Service.ServiceName: name,
				sfhttprequest.ExtraKey.HTTP.URL:            regConn.target,
			})
	}

	conn, err := grpc.Dial(regConn.target, regConn.dialOptions...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CloseConnections closes all registered connections
func CloseConnections() error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	var firstErr error
	for name, regConn := range globalRegistry.connections {
		if regConn.conn != nil {
			if err := regConn.conn.Close(); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("error closing connection '%s': %w", name, err)
			}
			regConn.conn = nil
		}
	}

	return firstErr
}

// CloseConnection closes a specific connection
func CloseConnection(name string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	regConn, exists := globalRegistry.connections[name]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", name)
	}

	if regConn.conn != nil {
		err := regConn.conn.Close()
		regConn.conn = nil
		return err
	}

	return nil
}

// FromService creates a new request builder for a specific service method using connection name, service name, and method name
func FromService(connectionName, serviceName, methodName string) *Request {
	// Create a default request
	r := &Request{
		ctx: context.Background(),
		md:  metadata.MD{},
	}

	// Get the connection
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[connectionName]
	autoConnectTimeout := globalRegistry.autoConnectTimeoutMs
	globalRegistry.mu.RUnlock()

	if !exists {
		r.err = fmt.Errorf("connection with name '%s' not registered", connectionName)
		return r
	}

	// If connection not established, try to wait for it if auto-connect is enabled
	globalRegistry.mu.RLock()
	grpcConn := conn.conn
	globalRegistry.mu.RUnlock()

	if grpcConn == nil && autoConnectTimeout > 0 {
		// Try to wait for the connection to establish
		if err := WaitForConnection(connectionName, autoConnectTimeout); err != nil {
			r.err = fmt.Errorf("connection '%s' is not established: %w", connectionName, err)
			return r
		}

		// Get the connection again after waiting
		globalRegistry.mu.RLock()
		conn = globalRegistry.connections[connectionName]
		grpcConn = conn.conn
		globalRegistry.mu.RUnlock()
	}

	// Lock again for service lookup
	globalRegistry.mu.RLock()

	// Find the service and method
	svcDef, serviceExists := conn.services[serviceName]
	if !serviceExists {
		globalRegistry.mu.RUnlock()
		r.err = fmt.Errorf("service '%s' not found in connection '%s'", serviceName, connectionName)
		return r
	}

	// Find the method in this service
	methodPath, methodExists := svcDef.Methods[methodName]
	if !methodExists {
		globalRegistry.mu.RUnlock()
		r.err = fmt.Errorf("method '%s' not found in service '%s'", methodName, serviceName)
		return r
	}

	// Get the physical connection
	if grpcConn == nil {
		globalRegistry.mu.RUnlock()
		target := conn.target
		r.err = fmt.Errorf("connection '%s' to target '%s' is not established, please check if the server is running",
			connectionName, target)
		return r
	}

	// Update the request with connection details
	r.conn = grpcConn
	r.methodName = methodPath
	r.logger = globalRegistry.logger

	// Add global headers to the request
	for key, values := range globalRegistry.globalHeaders {
		for _, value := range values {
			r.md.Append(key, value)
		}
	}

	// Add connection-specific headers
	for key, values := range conn.connectionHeaders {
		// Remove existing values from global headers that might exist
		if _, ok := r.md[key]; ok {
			delete(r.md, key)
		}
		// Add connection-specific values
		for _, value := range values {
			r.md.Append(key, value)
		}
	}
	globalRegistry.mu.RUnlock()

	return r
}

// FromServiceWithKey is for backward compatibility, uses the original two-argument signature
// but delegates to the new three-argument function
func FromServiceWithKey(connectionName, methodKey string) *Request {
	// Create a default request with error
	r := &Request{
		ctx: context.Background(),
		md:  metadata.MD{},
	}

	// Get the connection
	globalRegistry.mu.RLock()
	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		globalRegistry.mu.RUnlock()
		r.err = fmt.Errorf("connection with name '%s' not registered", connectionName)
		return r
	}

	// Find the service and method by method key
	var foundServiceName string
	var foundMethodName string

	// Search through all services in the connection to find the method key
	methodFound := false
	for serviceName, svcDef := range conn.services {
		for mName := range svcDef.Methods {
			if mName == methodKey {
				foundServiceName = serviceName
				foundMethodName = mName
				methodFound = true
				break
			}
		}
		if methodFound {
			break
		}
	}
	globalRegistry.mu.RUnlock()

	if !methodFound {
		r.err = fmt.Errorf("method key '%s' not found in connection '%s'", methodKey, connectionName)
		return r
	}

	// Delegate to the new function
	return FromService(connectionName, foundServiceName, foundMethodName)
}

// SetGlobalHeader sets a header that will be automatically added to all gRPC requests
func SetGlobalHeader(key, value string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.globalHeaders.Set(key, value)
}

// SetGlobalHeaders sets multiple headers that will be automatically added to all gRPC requests
func SetGlobalHeaders(headers map[string]string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	for key, value := range headers {
		globalRegistry.globalHeaders.Set(key, value)
	}
}

// GetGlobalHeaders returns all registered global headers
func GetGlobalHeaders() metadata.MD {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	// Create a copy of the global headers to avoid concurrent access issues
	headers := metadata.MD{}
	for key, values := range globalRegistry.globalHeaders {
		headers[key] = make([]string, len(values))
		copy(headers[key], values)
	}

	return headers
}

// ClearGlobalHeaders removes all global headers
func ClearGlobalHeaders() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.globalHeaders = metadata.MD{}
}

// RemoveGlobalHeader removes a specific global header
func RemoveGlobalHeader(key string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	delete(globalRegistry.globalHeaders, key)
}

// SetConnectionHeader sets a header that will be automatically added to all requests for a specific connection
func SetConnectionHeader(connectionName, key, value string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", connectionName)
	}

	conn.connectionHeaders.Set(key, value)
	return nil
}

// SetConnectionHeaders sets multiple headers that will be automatically added to all requests for a specific connection
func SetConnectionHeaders(connectionName string, headers map[string]string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", connectionName)
	}

	for key, value := range headers {
		conn.connectionHeaders.Set(key, value)
	}

	return nil
}

// Health checks if all gRPC connections are alive
func Health(ctx context.Context) error {
	globalRegistry.mu.RLock()
	connections := make([]string, 0, len(globalRegistry.connections))
	for name := range globalRegistry.connections {
		connections = append(connections, name)
	}
	globalRegistry.mu.RUnlock()

	if len(connections) == 0 {
		return fmt.Errorf("no gRPC connections registered")
	}

	var errs []string
	for _, name := range connections {
		globalRegistry.mu.RLock()
		conn, exists := globalRegistry.connections[name]
		grpcConn := conn.conn
		globalRegistry.mu.RUnlock()

		if !exists || grpcConn == nil {
			errs = append(errs, fmt.Sprintf("connection '%s' is not established", name))
			continue
		}

		// Check connection state
		state := grpcConn.GetState()
		if state == connectivity.Shutdown || state == connectivity.TransientFailure {
			errs = append(errs, fmt.Sprintf("connection '%s' is in bad state: %v", name, state))
			continue
		}

		// Optional: force connectivity check with a short deadline
		checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		if !grpcConn.WaitForStateChange(checkCtx, state) {
			errs = append(errs, fmt.Sprintf("connection '%s' connectivity check timed out", name))
		}
		cancel()
	}

	if len(errs) > 0 {
		return fmt.Errorf("gRPC health check failed: %s", strings.Join(errs, "; "))
	}

	return nil
}

// GetService returns a service client by name
func GetService[T any](serviceName string) (T, error) {
	var zeroVal T

	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	client, exists := globalRegistry.serviceClients[serviceName]
	if !exists {
		return zeroVal, fmt.Errorf("service '%s' not registered", serviceName)
	}

	if typedClient, ok := client.(T); ok {
		return typedClient, nil
	}

	return zeroVal, fmt.Errorf("service '%s' is not of the requested type", serviceName)
}

// RegisterService registers a service with a connection
func RegisterService(serviceName, connectionName string, clientConstructor interface{}) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if globalRegistry.serviceClients == nil {
		globalRegistry.serviceClients = make(map[string]interface{})
	}

	// Create client and store it
	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		return fmt.Errorf("connection '%s' not registered", connectionName)
	}

	if conn.conn == nil {
		return fmt.Errorf("connection '%s' is not established", connectionName)
	}

	// Create the client using reflection
	constructorVal := reflect.ValueOf(clientConstructor)
	args := []reflect.Value{reflect.ValueOf(conn.conn)}
	results := constructorVal.Call(args)

	if len(results) != 1 {
		return fmt.Errorf("client constructor should return exactly one value")
	}

	client := results[0].Interface()
	globalRegistry.serviceClients[serviceName] = client

	return nil
}

// WaitForConnection waits for a specific connection to be established with a timeout
func WaitForConnection(connectionName string, timeoutMs int) error {
	deadline := time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)

	for time.Now().Before(deadline) {
		globalRegistry.mu.RLock()
		conn, exists := globalRegistry.connections[connectionName]
		established := exists && conn.conn != nil
		globalRegistry.mu.RUnlock()

		if established {
			return nil
		}

		// Check every 100ms
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for connection '%s' to be established", connectionName)
}
