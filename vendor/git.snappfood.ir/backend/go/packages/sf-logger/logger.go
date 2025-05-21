package sflogger

import (
	"context"
)

// Logger defines the core logging interface
type Logger interface {
	// Core logging methods
	Debug(msg string, extra map[string]interface{})
	Info(msg string, extra map[string]interface{})
	Warn(msg string, extra map[string]interface{})
	Error(msg string, extra map[string]interface{})
	Fatal(msg string, extra map[string]interface{})

	// Formatted logging methods
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	// Context-aware logging methods
	DebugContext(ctx context.Context, msg string, extra map[string]interface{})
	InfoContext(ctx context.Context, msg string, extra map[string]interface{})
	WarnContext(ctx context.Context, msg string, extra map[string]interface{})
	ErrorContext(ctx context.Context, msg string, extra map[string]interface{})
	FatalContext(ctx context.Context, msg string, extra map[string]interface{})

	// Category-based logging methods
	DebugWithCategory(cat string, sub string, msg string, extra map[string]interface{})
	InfoWithCategory(cat string, sub string, msg string, extra map[string]interface{})
	WarnWithCategory(cat string, sub string, msg string, extra map[string]interface{})
	ErrorWithCategory(cat string, sub string, msg string, extra map[string]interface{})
	FatalWithCategory(cat string, sub string, msg string, extra map[string]interface{})
}

// Categories provides structured access to all log categories
type Categories struct {
	// System category group
	System struct {
		General       string
		Internal      string
		Startup       string
		Shutdown      string
		Health        string
		Configuration string
	}

	// Infrastructure category group
	Infrastructure struct {
		IO             string
		Network        string
		Infrastructure string
		Cloud          string
		Kubernetes     string
		Container      string
	}

	// Database category group
	Database struct {
		Database string
		MySQL    string
		MongoDB  string
		Postgres string
		Redis    string
		Cache    string
	}

	// API category group
	API struct {
		API           string
		HTTP          string
		GRPC          string
		REST          string
		GraphQL       string
		WebSocket     string
		Messaging     string
		Documentation string
	}

	// Security category group
	Security struct {
		Security       string
		Authentication string
		Authorization  string
	}

	// Business category group
	Business struct {
		Business    string
		Validation  string
		Transaction string
		Workflow    string
		Payment     string
		Order       string
		User        string
	}

	// Monitoring category group
	Monitoring struct {
		Metrics         string
		Traces          string
		Prometheus      string
		RequestResponse string
		Performance     string
	}

	// Error category group
	Error struct {
		Error     string
		Exception string
		Retry     string
		Circuit   string
		Fallback  string
		RateLimit string
	}

	// Service category group
	Service struct {
		Service     string
		External    string
		Dependency  string
		Integration string
	}

	// Development category group
	Development struct {
		Development string
		Testing     string
		Debugging   string
	}
}

// SubCategories provides structured access to all log subcategories
type SubCategories struct {
	// Operation subcategory group
	Operation struct {
		Startup        string
		Shutdown       string
		Initialization string
		Configuration  string
		Registration   string
		Discovery      string
		Provisioning   string
	}

	// Status subcategory group
	Status struct {
		Success   string
		Failure   string
		Warning   string
		Info      string
		Debug     string
		Error     string
		Timeout   string
		Retry     string
		Reconnect string
	}

	// Database subcategory group
	Database struct {
		Query       string
		Transaction string
		Migration   string
		Backup      string
		Restore     string
		Select      string
		Insert      string
		Update      string
		Delete      string
		Rollback    string
	}

	// API subcategory group
	API struct {
		Request         string
		Response        string
		Validation      string
		Parsing         string
		Serialization   string
		Deserialization string
	}

	// Security subcategory group
	Security struct {
		Authentication string
		Authorization  string
		Token          string
		Session        string
		Encryption     string
		Decryption     string
	}

	// Networking subcategory group
	Networking struct {
		Connection    string
		Disconnection string
		Latency       string
		Throughput    string
		Scaling       string
		ResourceUsage string
	}

	// Resilience subcategory group
	Resilience struct {
		CircuitBreaker string
		Fallback       string
		RateLimit      string
		Throttling     string
		Bulkhead       string
	}

	// User subcategory group
	User struct {
		Creation string
		Deletion string
		Update   string
		Login    string
		Logout   string
		Activity string
	}

	// Business subcategory group
	Business struct {
		OrderCreation    string
		OrderFulfillment string
		Payment          string
		Shipping         string
		Notification     string
		Workflow         string
	}

	// Integration subcategory group
	Integration struct {
		APICall     string
		ServiceCall string
		Webhook     string
		Integration string
		ThirdParty  string
	}

	// Legacy subcategory group for backward compatibility
	Legacy struct {
		Api                 string
		ServerShutdown      string
		MobileValidation    string
		DefaultRoleNotFound string
	}
}

// ExtraKeys provides structured access to all log extra keys
type ExtraKeys struct {
	// Metadata keys group
	Metadata struct {
		AppName     string
		LoggerName  string
		Version     string
		Environment string
		PodName     string
		NodeName    string
		InstanceID  string
	}

	// Network keys group
	Network struct {
		ClientIP string
		HostIP   string
		Port     string
		Protocol string
	}

	// Request keys group
	Request struct {
		RequestID     string
		TraceID       string
		SpanID        string
		UserID        string
		SessionID     string
		CorrelationID string
	}

	// HTTP keys group
	HTTP struct {
		Method        string
		StatusCode    string
		Path          string
		URL           string
		UserAgent     string
		Referer       string
		ContentType   string
		ContentLength string
	}

	// Performance keys group
	Performance struct {
		Latency  string
		Duration string
		BodySize string
	}

	// Payload keys group
	Payload struct {
		RequestBody  string
		ResponseBody string
		Payload      string
	}

	// Error keys group
	Error struct {
		ErrorMessage string
		ErrorCode    string
		StackTrace   string
		Exception    string
	}

	// Database keys group
	Database struct {
		Query         string
		Table         string
		TransactionID string
		RowsAffected  string
		ResultCount   string
	}

	// Service keys group
	Service struct {
		ServiceName    string
		ServiceVersion string
		Endpoint       string
	}

	// Business keys group
	Business struct {
		OrderID    string
		CustomerID string
		ProductID  string
		AccountID  string
	}

	// Resource keys group
	Resource struct {
		CPUUsage        string
		MemoryUsage     string
		DiskUsage       string
		ConnectionCount string
	}

	// Legacy keys group for backward compatibility
	Legacy struct {
		AppName      string
		LoggerName   string
		ClientIp     string
		HostIp       string
		Method       string
		StatusCode   string
		BodySize     string
		Path         string
		Latency      string
		RequestBody  string
		ResponseBody string
		ErrorMessage string
	}
}

// Category is a global variable for accessing all log categories
var Category = initCategories()

// SubCategory is a global variable for accessing all log subcategories
var SubCategory = initSubCategories()

// ExtraKey is a global variable for accessing all log extra keys
var ExtraKey = initExtraKeys()

// Initialize Categories struct with values
func initCategories() Categories {
	c := Categories{}

	// System
	c.System.General = "General"
	c.System.Internal = "Internal"
	c.System.Startup = "Startup"
	c.System.Shutdown = "Shutdown"
	c.System.Health = "Health"
	c.System.Configuration = "Configuration"

	// Infrastructure
	c.Infrastructure.IO = "IO"
	c.Infrastructure.Network = "Network"
	c.Infrastructure.Infrastructure = "Infrastructure"
	c.Infrastructure.Cloud = "Cloud"
	c.Infrastructure.Kubernetes = "Kubernetes"
	c.Infrastructure.Container = "Container"

	// Database
	c.Database.Database = "Database"
	c.Database.MySQL = "Mysql"
	c.Database.MongoDB = "MongoDB"
	c.Database.Postgres = "Postgres"
	c.Database.Redis = "Redis"
	c.Database.Cache = "Cache"

	// API
	c.API.API = "API"
	c.API.HTTP = "HTTP"
	c.API.GRPC = "GRPC"
	c.API.REST = "REST"
	c.API.GraphQL = "GraphQL"
	c.API.WebSocket = "WebSocket"
	c.API.Messaging = "Messaging"
	c.API.Documentation = "Documentation"

	// Security
	c.Security.Security = "Security"
	c.Security.Authentication = "Authentication"
	c.Security.Authorization = "Authorization"

	// Business
	c.Business.Business = "Business"
	c.Business.Validation = "Validation"
	c.Business.Transaction = "Transaction"
	c.Business.Workflow = "Workflow"
	c.Business.Payment = "Payment"
	c.Business.Order = "Order"
	c.Business.User = "User"

	// Monitoring
	c.Monitoring.Metrics = "Metrics"
	c.Monitoring.Traces = "Traces"
	c.Monitoring.Prometheus = "Prometheus"
	c.Monitoring.RequestResponse = "RequestResponse"
	c.Monitoring.Performance = "Performance"

	// Error
	c.Error.Error = "Error"
	c.Error.Exception = "Exception"
	c.Error.Retry = "Retry"
	c.Error.Circuit = "Circuit"
	c.Error.Fallback = "Fallback"
	c.Error.RateLimit = "RateLimit"

	// Service
	c.Service.Service = "Service"
	c.Service.External = "ExternalService"
	c.Service.Dependency = "Dependency"
	c.Service.Integration = "Integration"

	// Development
	c.Development.Development = "Development"
	c.Development.Testing = "Testing"
	c.Development.Debugging = "Debugging"

	return c
}

// Initialize SubCategories struct with values
func initSubCategories() SubCategories {
	s := SubCategories{}

	// Operation
	s.Operation.Startup = "Startup"
	s.Operation.Shutdown = "Shutdown"
	s.Operation.Initialization = "Initialization"
	s.Operation.Configuration = "Configuration"
	s.Operation.Registration = "Registration"
	s.Operation.Discovery = "Discovery"
	s.Operation.Provisioning = "Provisioning"

	// Status
	s.Status.Success = "Success"
	s.Status.Failure = "Failure"
	s.Status.Warning = "Warning"
	s.Status.Info = "Info"
	s.Status.Debug = "Debug"
	s.Status.Error = "Error"
	s.Status.Timeout = "Timeout"
	s.Status.Retry = "Retry"
	s.Status.Reconnect = "Reconnect"

	// Database
	s.Database.Query = "Query"
	s.Database.Transaction = "Transaction"
	s.Database.Migration = "Migration"
	s.Database.Backup = "Backup"
	s.Database.Restore = "Restore"
	s.Database.Select = "Select"
	s.Database.Insert = "Insert"
	s.Database.Update = "Update"
	s.Database.Delete = "Delete"
	s.Database.Rollback = "Rollback"

	// API
	s.API.Request = "Request"
	s.API.Response = "Response"
	s.API.Validation = "Validation"
	s.API.Parsing = "Parsing"
	s.API.Serialization = "Serialization"
	s.API.Deserialization = "Deserialization"

	// Security
	s.Security.Authentication = "Authentication"
	s.Security.Authorization = "Authorization"
	s.Security.Token = "Token"
	s.Security.Session = "Session"
	s.Security.Encryption = "Encryption"
	s.Security.Decryption = "Decryption"

	// Networking
	s.Networking.Connection = "Connection"
	s.Networking.Disconnection = "Disconnection"
	s.Networking.Latency = "Latency"
	s.Networking.Throughput = "Throughput"
	s.Networking.Scaling = "Scaling"
	s.Networking.ResourceUsage = "ResourceUsage"

	// Resilience
	s.Resilience.CircuitBreaker = "CircuitBreaker"
	s.Resilience.Fallback = "Fallback"
	s.Resilience.RateLimit = "RateLimit"
	s.Resilience.Throttling = "Throttling"
	s.Resilience.Bulkhead = "Bulkhead"

	// User
	s.User.Creation = "UserCreation"
	s.User.Deletion = "UserDeletion"
	s.User.Update = "UserUpdate"
	s.User.Login = "UserLogin"
	s.User.Logout = "UserLogout"
	s.User.Activity = "UserActivity"

	// Business
	s.Business.OrderCreation = "OrderCreation"
	s.Business.OrderFulfillment = "OrderFulfillment"
	s.Business.Payment = "Payment"
	s.Business.Shipping = "Shipping"
	s.Business.Notification = "Notification"
	s.Business.Workflow = "Workflow"

	// Integration
	s.Integration.APICall = "APICall"
	s.Integration.ServiceCall = "ServiceCall"
	s.Integration.Webhook = "Webhook"
	s.Integration.Integration = "Integration"
	s.Integration.ThirdParty = "ThirdParty"

	// Legacy
	s.Legacy.Api = "Api"
	s.Legacy.ServerShutdown = "ServerShutdown"
	s.Legacy.MobileValidation = "MobileValidation"
	s.Legacy.DefaultRoleNotFound = "DefaultRoleNotFound"

	return s
}

// Initialize ExtraKeys struct with values
func initExtraKeys() ExtraKeys {
	k := ExtraKeys{}

	// Metadata
	k.Metadata.AppName = "AppName"
	k.Metadata.LoggerName = "Logger"
	k.Metadata.Version = "Version"
	k.Metadata.Environment = "Environment"
	k.Metadata.PodName = "PodName"
	k.Metadata.NodeName = "NodeName"
	k.Metadata.InstanceID = "InstanceID"

	// Network
	k.Network.ClientIP = "ClientIP"
	k.Network.HostIP = "HostIP"
	k.Network.Port = "Port"
	k.Network.Protocol = "Protocol"

	// Request
	k.Request.RequestID = "RequestID"
	k.Request.TraceID = "TraceID"
	k.Request.SpanID = "SpanID"
	k.Request.UserID = "UserID"
	k.Request.SessionID = "SessionID"
	k.Request.CorrelationID = "CorrelationID"

	// HTTP
	k.HTTP.Method = "Method"
	k.HTTP.StatusCode = "StatusCode"
	k.HTTP.Path = "Path"
	k.HTTP.URL = "URL"
	k.HTTP.UserAgent = "UserAgent"
	k.HTTP.Referer = "Referer"
	k.HTTP.ContentType = "ContentType"
	k.HTTP.ContentLength = "ContentLength"

	// Performance
	k.Performance.Latency = "Latency"
	k.Performance.Duration = "Duration"
	k.Performance.BodySize = "BodySize"

	// Payload
	k.Payload.RequestBody = "RequestBody"
	k.Payload.ResponseBody = "ResponseBody"
	k.Payload.Payload = "Payload"

	// Error
	k.Error.ErrorMessage = "ErrorMessage"
	k.Error.ErrorCode = "ErrorCode"
	k.Error.StackTrace = "StackTrace"
	k.Error.Exception = "Exception"

	// Database
	k.Database.Query = "Query"
	k.Database.Table = "Table"
	k.Database.TransactionID = "TransactionID"
	k.Database.RowsAffected = "RowsAffected"
	k.Database.ResultCount = "ResultCount"

	// Service
	k.Service.ServiceName = "ServiceName"
	k.Service.ServiceVersion = "ServiceVersion"
	k.Service.Endpoint = "Endpoint"

	// Business
	k.Business.OrderID = "OrderID"
	k.Business.CustomerID = "CustomerID"
	k.Business.ProductID = "ProductID"
	k.Business.AccountID = "AccountID"

	// Resource
	k.Resource.CPUUsage = "CPUUsage"
	k.Resource.MemoryUsage = "MemoryUsage"
	k.Resource.DiskUsage = "DiskUsage"
	k.Resource.ConnectionCount = "ConnectionCount"

	// Legacy
	k.Legacy.AppName = "AppName"
	k.Legacy.LoggerName = "Logger"
	k.Legacy.ClientIp = "ClientIp"
	k.Legacy.HostIp = "HostIp"
	k.Legacy.Method = "Method"
	k.Legacy.StatusCode = "StatusCode"
	k.Legacy.BodySize = "BodySize"
	k.Legacy.Path = "Path"
	k.Legacy.Latency = "Latency"
	k.Legacy.RequestBody = "RequestBody"
	k.Legacy.ResponseBody = "ResponseBody"
	k.Legacy.ErrorMessage = "ErrorMessage"

	return k
}