package httpo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	sfhttprequest "git.snappfood.ir/backend/go/packages/sf-http-request"
	"go.elastic.co/apm/module/apmhttp/v2"
)

// Registry is a global registry for HTTP connections
type Registry struct {
	mu               sync.RWMutex
	connections      map[string]*registeredConnection
	defaultTransport http.RoundTripper
	globalHeaders    map[string]string
	logger           sfhttprequest.Logger
}

// Option is a function to customize the Registry
type Option func(*Registry)

// WithDefaultTransport sets the default transport for all new connections
func WithDefaultTransport(transport http.RoundTripper) Option {
	return func(r *Registry) {
		r.defaultTransport = transport
	}
}

// WithGlobalHeaders sets the global headers for all connections
func WithGlobalHeaders(headers map[string]string) Option {
	return func(r *Registry) {
		r.globalHeaders = headers
	}
}

// WithConnectionDetails sets the connection details for a named connection
func WithConnectionDetails(name, baseURL string, options ...ClientOption) Option {
	return func(r *Registry) {
		if r.connections == nil {
			r.connections = make(map[string]*registeredConnection)
		}

		// Validate URL
		_, err := url.Parse(baseURL)
		if err != nil {
			// Just log the error if logger is available, but still register the connection
			if r.logger != nil {
				r.logger.ErrorWithCategory(
					sfhttprequest.Category.API.HTTP,
					sfhttprequest.SubCategory.Status.Error,
					"Invalid URL format",
					map[string]interface{}{
						sfhttprequest.ExtraKey.Service.ServiceName: name,
						sfhttprequest.ExtraKey.HTTP.URL:            baseURL,
						sfhttprequest.ExtraKey.Error.ErrorMessage:  err.Error(),
					})
			}
		}

		r.connections[name] = &registeredConnection{
			baseURL:           baseURL,
			options:           options,
			connectionHeaders: make(map[string]string),
		}
	}
}

// WithLogger sets the logger for the registry
func WithLogger(logger sfhttprequest.Logger) Option {
	return func(r *Registry) {
		r.logger = logger
	}
}

// RegisterConnection configures the HTTP connections with the provided options
func RegisterConnection(opts ...Option) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// Apply options
	for _, opt := range opts {
		opt(globalRegistry)
	}

	if globalRegistry.logger != nil {
		globalRegistry.logger.InfoWithCategory(
			sfhttprequest.Category.API.HTTP,
			sfhttprequest.SubCategory.Operation.Registration,
			"HTTP registry configured",
			map[string]interface{}{
				sfhttprequest.ExtraKey.Service.ServiceName: "HTTP registry",
			})
	}

	return nil
}

type registeredConnection struct {
	client            *http.Client
	baseURL           string
	options           []ClientOption
	connectionHeaders map[string]string
}

// Global registry instance
var globalRegistry = &Registry{
	connections:      make(map[string]*registeredConnection),
	defaultTransport: http.DefaultTransport,
	globalHeaders:    make(map[string]string),
}

// SetDefaultTransport sets the default transport option for all new connections
func SetDefaultTransport(transport http.RoundTripper) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.defaultTransport = transport
}

// getConnection returns the connection with the given name, initializing it if needed
func getConnection(name string) (*http.Client, string, error) {
	globalRegistry.mu.RLock()
	regConn, exists := globalRegistry.connections[name]
	logger := globalRegistry.logger
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, "", fmt.Errorf("connection with name '%s' not registered", name)
	}

	// If client is already initialized, return it
	if regConn.client != nil {
		return regConn.client, regConn.baseURL, nil
	}

	// Initialize the client
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	// Check again in case another goroutine initialized it
	if regConn.client != nil {
		return regConn.client, regConn.baseURL, nil
	}

	// Create a transport with APM tracing
	transport := globalRegistry.defaultTransport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Create a new client with the transport
	client := &http.Client{
		Transport: transport,
	}

	// Apply custom options to the client
	for _, option := range regConn.options {
		option(client)
	}

	// Always wrap the client with APM tracing, regardless of other options
	client = apmhttp.WrapClient(client)

	regConn.client = client
	globalRegistry.connections[name] = regConn

	if logger != nil {
		logger.InfoWithCategory(
			sfhttprequest.Category.API.HTTP,
			sfhttprequest.SubCategory.Networking.Connection,
			"HTTP client initialized",
			map[string]interface{}{
				sfhttprequest.ExtraKey.Service.ServiceName: name,
				sfhttprequest.ExtraKey.HTTP.URL:            regConn.baseURL,
			})
	}

	return client, regConn.baseURL, nil
}

// CloseConnections closes all registered connections
// Note: In HTTP there's no explicit connection closing like in gRPC,
// but this method can be used to clean up resources
func CloseConnections() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	for name, regConn := range globalRegistry.connections {
		regConn.client = nil
		globalRegistry.connections[name] = regConn
	}
}

// CloseConnection closes a specific connection by name
func CloseConnection(name string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	regConn, exists := globalRegistry.connections[name]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", name)
	}

	regConn.client = nil
	globalRegistry.connections[name] = regConn

	return nil
}

// SetGlobalHeader sets a header that will be automatically added to all requests
func SetGlobalHeader(key, value string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.globalHeaders[key] = value
}

// SetGlobalHeaders sets multiple headers that will be automatically added to all requests
func SetGlobalHeaders(headers map[string]string) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	for key, value := range headers {
		globalRegistry.globalHeaders[key] = value
	}
}

// GetGlobalHeaders returns all registered global headers
func GetGlobalHeaders() map[string]string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	// Create a copy of the global headers map to avoid concurrent access issues
	headers := make(map[string]string, len(globalRegistry.globalHeaders))
	for key, value := range globalRegistry.globalHeaders {
		headers[key] = value
	}

	return headers
}

// ClearGlobalHeaders removes all global headers
func ClearGlobalHeaders() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.globalHeaders = make(map[string]string)
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

	if conn.connectionHeaders == nil {
		conn.connectionHeaders = make(map[string]string)
	}

	conn.connectionHeaders[key] = value
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

	if conn.connectionHeaders == nil {
		conn.connectionHeaders = make(map[string]string)
	}

	for key, value := range headers {
		conn.connectionHeaders[key] = value
	}

	return nil
}

// GetConnectionHeaders returns all registered headers for a specific connection
func GetConnectionHeaders(connectionName string) (map[string]string, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		return nil, fmt.Errorf("connection with name '%s' not registered", connectionName)
	}

	// Create a copy of the headers map to avoid concurrent access issues
	headers := make(map[string]string, len(conn.connectionHeaders))
	for key, value := range conn.connectionHeaders {
		headers[key] = value
	}

	return headers, nil
}

// ClearConnectionHeaders removes all headers for a specific connection
func ClearConnectionHeaders(connectionName string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", connectionName)
	}

	conn.connectionHeaders = make(map[string]string)
	return nil
}

// RemoveConnectionHeader removes a specific header for a connection
func RemoveConnectionHeader(connectionName, key string) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	conn, exists := globalRegistry.connections[connectionName]
	if !exists {
		return fmt.Errorf("connection with name '%s' not registered", connectionName)
	}

	delete(conn.connectionHeaders, key)
	return nil
}

// FromConnection creates a new request builder using a registered connection identified by name
func FromConnection(name string) (*Request, error) {
	client, baseURL, err := getConnection(name)
	if err != nil {
		return nil, err
	}

	r := &Request{
		client:      client,
		url:         baseURL,
		headers:     make(map[string]string),
		queryParams: url.Values{},
		formData:    url.Values{},
		ctx:         context.Background(),
	}

	// Add global headers to the request
	globalRegistry.mu.RLock()
	for key, value := range globalRegistry.globalHeaders {
		r.headers[key] = value
	}

	// Add connection-specific headers (these override globals with the same key)
	conn, exists := globalRegistry.connections[name]
	if exists {
		// Set the logger from the registry
		r.logger = globalRegistry.logger

		if conn.connectionHeaders != nil {
			for key, value := range conn.connectionHeaders {
				r.headers[key] = value
			}
		}
	}
	globalRegistry.mu.RUnlock()

	return r, nil
}

// Health checks if all HTTP connections are alive
func Health(ctx context.Context) error {
	globalRegistry.mu.RLock()
	connections := make([]string, 0, len(globalRegistry.connections))
	for name := range globalRegistry.connections {
		connections = append(connections, name)
	}
	globalRegistry.mu.RUnlock()

	if len(connections) == 0 {
		return fmt.Errorf("no HTTP connections registered")
	}

	var errs []string
	for _, name := range connections {
		client, baseURL, err := getConnection(name)
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' error: %v", name, err))
			continue
		}

		// Create a simple GET request to check if the connection is alive
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' request creation error: %v", name, err))
			continue
		}

		// Send a HEAD request to the base URL with a timeout
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, fmt.Sprintf("connection '%s' request error: %v", name, err))
			continue
		}
		resp.Body.Close() // Always close the response body

		// Check for non-successful status codes (this can be adjusted based on your needs)
		if resp.StatusCode >= 500 {
			errs = append(errs, fmt.Sprintf("connection '%s' returned status %d", name, resp.StatusCode))
			continue
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("HTTP health check failed: %s", strings.Join(errs, "; "))
	}

	return nil
}
