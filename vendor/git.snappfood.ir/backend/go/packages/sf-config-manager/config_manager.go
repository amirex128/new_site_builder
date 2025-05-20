package sfconfigmanager

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

// ManagerType defines the type of configuration manager
type ManagerType string

// Available manager types
const (
	VaultManager     ManagerType = "vault"
	EnvManager       ManagerType = "env"
	FileManager      ManagerType = "file"
	ConsulManager    ManagerType = "consul"
	EtcdManager      ManagerType = "etcd"
	ZookeeperManager ManagerType = "zookeeper"
)

// config holds common configuration for configuration manager
type config struct {
	// Primary config source (e.g., "vault", "env", "file")
	Type ManagerType

	// Options specific to the implementation
	Options map[string]interface{}
}

// connectionDetails contains the necessary details to connect to a config source
type connectionDetails struct {
	Type     ManagerType
	Address  string
	Username string
	Password string
	Options  map[string]interface{}
}

// configLoader represents the main configuration loader
type configLoader struct {
	managers     []configManager
	log          Logger
	config       interface{}
	retryOptions *RetryOptions
}

// configManager defines the interface for configuration manager operations
type configManager interface {
	// Load loads configuration into the provided struct
	Load(target interface{}) error

	// Get gets a specific configuration value by key
	Get(key string) (interface{}, error)

	// GetString gets a specific configuration value as string
	GetString(key string) (string, error)

	// GetInt gets a specific configuration value as int
	GetInt(key string) (int, error)

	// GetBool gets a specific configuration value as bool
	GetBool(key string) (bool, error)

	// Close cleans up resources used by the manager
	Close() error
}

// provider manages the current active configuration manager
type provider struct {
	manager configManager
	mu      sync.RWMutex
}

// Global provider instance
var configProvider = &provider{}

// RegisterConnection initializes a new configuration manager with the provided options
func RegisterConnection(options ...ConfigOption) (interface{}, error) {
	loader := &configLoader{
		managers: []configManager{},
	}

	// Apply all provided options
	for _, option := range options {
		option(loader)
	}

	// If no config was provided, return an error
	if loader.config == nil {
		return nil, fmt.Errorf("no configuration struct provided")
	}

	// If no retry options were provided, use defaults
	if loader.retryOptions == nil {
		loader.retryOptions = DefaultRetryOptions()
	}

	// Try to load from each manager in order
	var lastErr error
	for _, manager := range loader.managers {
		// Create a load operation wrapped with retry logic
		operationName := fmt.Sprintf("load configuration from %T", manager)

		err := WithRetry(
			func() error {
				return manager.Load(loader.config)
			},
			loader.retryOptions,
			loader.log,
			operationName,
		)

		if err != nil {
			lastErr = err
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, "Failed to load from manager after retries", extraMap)
			}
			continue
		}

		// Successfully loaded, set this manager as the current one
		configProvider.mu.Lock()
		configProvider.manager = manager
		configProvider.mu.Unlock()

		return loader.config, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", lastErr)
	}

	return nil, fmt.Errorf("no valid configuration managers could be initialized")
}

// ConfigOption represents a functional option for configuring the ConfigLoader
type ConfigOption func(*configLoader)

// WithConfig sets the configuration struct to be populated
func WithConfig(config interface{}) ConfigOption {
	return func(loader *configLoader) {
		loader.config = config
	}
}

// WithLogger sets the logger to be used by all managers
func WithLogger(logger Logger) ConfigOption {
	return func(loader *configLoader) {
		loader.log = logger
	}
}

// VaultOptions contains Vault-specific configuration options
type VaultOptions struct {
	Address     string
	Token       string
	SecretPath  string
	SecretMount string
	Timeout     time.Duration
	AuthMethod  string
	Role        string
}

// WithVaultOptions adds Vault-specific options
func WithVaultOptions(address string, token string, options *VaultOptions) ConfigOption {
	return func(loader *configLoader) {
		// Create manager-specific config
		cfg := config{
			Type:    VaultManager,
			Options: make(map[string]interface{}),
		}

		vaultOptions := map[string]interface{}{
			"address": address,
			"token":   token,
		}

		if options != nil {
			if options.SecretPath != "" {
				vaultOptions["secretPath"] = options.SecretPath
			}
			if options.SecretMount != "" {
				vaultOptions["secretMount"] = options.SecretMount
			}
			if options.Timeout > 0 {
				vaultOptions["timeout"] = options.Timeout
			}
			if options.AuthMethod != "" {
				vaultOptions["authMethod"] = options.AuthMethod
			}
			if options.Role != "" {
				vaultOptions["role"] = options.Role
			}
		}

		cfg.Options["vault"] = vaultOptions

		// Initialize manager
		manager, err := NewVaultManager(Config(cfg), loader.log)
		if err != nil {
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to initialize Vault manager"), extraMap)
			}
			return
		}

		loader.managers = append(loader.managers, manager)
	}
}

// FileOptions contains File-specific configuration options
type FileOptions struct {
	Path string
	Type string // json, yaml
}

// WithFileOptions adds File-specific options
func WithFileOptions(path string, options *FileOptions) ConfigOption {
	return func(loader *configLoader) {
		// Create manager-specific config
		cfg := config{
			Type:    FileManager,
			Options: make(map[string]interface{}),
		}

		fileOptions := map[string]interface{}{
			"path": path,
		}

		if options != nil && options.Type != "" {
			fileOptions["type"] = options.Type
		}

		cfg.Options["file"] = fileOptions

		// Initialize manager
		manager, err := NewFileManager(Config(cfg), loader.log)
		if err != nil {
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to initialize File manager"), extraMap)
			}
			return
		}

		loader.managers = append(loader.managers, manager)
	}
}

// EnvOptions contains Environment-specific configuration options
type EnvOptions struct {
	Prefix string
}

// WithEnvOptions adds Environment-specific options
func WithEnvOptions(options *EnvOptions) ConfigOption {
	return func(loader *configLoader) {
		// Create manager-specific config
		cfg := config{
			Type:    EnvManager,
			Options: make(map[string]interface{}),
		}

		envOptions := map[string]interface{}{}

		if options != nil && options.Prefix != "" {
			envOptions["prefix"] = options.Prefix
		}

		cfg.Options["env"] = envOptions

		// Initialize manager
		manager, err := NewEnvManager(Config(cfg), loader.log)
		if err != nil {
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to initialize Environment manager"), extraMap)
			}
			return
		}

		loader.managers = append(loader.managers, manager)
	}
}

// ConsulOptions contains Consul-specific configuration options
type ConsulOptions struct {
	Address    string
	Datacenter string
	Token      string
	Timeout    time.Duration
	Prefix     string
	Username   string
	Password   string
}

// WithConsulOptions adds Consul-specific options
func WithConsulOptions(address string, options *ConsulOptions) ConfigOption {
	return func(loader *configLoader) {
		// Create manager-specific config
		cfg := config{
			Type:    ConsulManager,
			Options: make(map[string]interface{}),
		}

		consulOptions := map[string]interface{}{
			"address": address,
		}

		if options != nil {
			if options.Datacenter != "" {
				consulOptions["datacenter"] = options.Datacenter
			}
			if options.Token != "" {
				consulOptions["token"] = options.Token
			}
			if options.Prefix != "" {
				consulOptions["prefix"] = options.Prefix
			}
			if options.Timeout > 0 {
				consulOptions["timeout"] = options.Timeout
			}
			if options.Username != "" {
				consulOptions["username"] = options.Username
			}
			if options.Password != "" {
				consulOptions["password"] = options.Password
			}
		}

		cfg.Options["consul"] = consulOptions

		// Initialize manager
		manager, err := NewConsulManager(Config(cfg), loader.log)
		if err != nil {
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to initialize Consul manager"), extraMap)
			}
			return
		}

		loader.managers = append(loader.managers, manager)
	}
}

// EtcdOptions contains etcd-specific configuration options
type EtcdOptions struct {
	Endpoints []string
	Username  string
	Password  string
	Timeout   time.Duration
	BasePath  string
}

// WithEtcdOptions adds etcd-specific options
func WithEtcdOptions(endpoints []string, options *EtcdOptions) ConfigOption {
	return func(loader *configLoader) {
		// Create manager-specific config
		cfg := config{
			Type:    EtcdManager,
			Options: make(map[string]interface{}),
		}

		etcdOptions := map[string]interface{}{
			"endpoints": endpoints,
		}

		if options != nil {
			if options.Username != "" {
				etcdOptions["username"] = options.Username
			}
			if options.Password != "" {
				etcdOptions["password"] = options.Password
			}
			if options.BasePath != "" {
				etcdOptions["prefix"] = options.BasePath
			}
			if options.Timeout > 0 {
				etcdOptions["timeout"] = options.Timeout
			}
		}

		cfg.Options["etcd"] = etcdOptions

		// Initialize manager
		manager, err := NewEtcdManager(Config(cfg), loader.log)
		if err != nil {
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to initialize etcd manager"), extraMap)
			}
			return
		}

		loader.managers = append(loader.managers, manager)
	}
}

// ZookeeperOptions contains Zookeeper-specific configuration options
type ZookeeperOptions struct {
	Servers  []string
	Timeout  time.Duration
	BasePath string
	Username string
	Password string
}

// WithZookeeperOptions adds Zookeeper-specific options
func WithZookeeperOptions(servers []string, options *ZookeeperOptions) ConfigOption {
	return func(loader *configLoader) {
		// Create manager-specific config
		cfg := config{
			Type:    ZookeeperManager,
			Options: make(map[string]interface{}),
		}

		zkOptions := map[string]interface{}{
			"servers": servers,
		}

		if options != nil {
			if options.Username != "" {
				zkOptions["username"] = options.Username
			}
			if options.Password != "" {
				zkOptions["password"] = options.Password
			}
			if options.BasePath != "" {
				zkOptions["basePath"] = options.BasePath
			}
			if options.Timeout > 0 {
				zkOptions["sessionTimeout"] = options.Timeout
			}
		}

		cfg.Options["zookeeper"] = zkOptions

		// Initialize manager
		manager, err := NewZookeeperManager(Config(cfg), loader.log)
		if err != nil {
			if loader.log != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				loader.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to initialize Zookeeper manager"), extraMap)
			}
			return
		}

		loader.managers = append(loader.managers, manager)
	}
}

// Backward compatibility functions for migrating from the old API

// WithConnectionDetails adds a new connection configuration (legacy compatibility)
func WithConnectionDetails(managerType ManagerType, address string, username string, password string, options ...interface{}) ConfigOption {
	return func(loader *configLoader) {
		switch managerType {
		case VaultManager:
			vaultOpts := &VaultOptions{
				Token: password,
			}
			if len(options) > 0 {
				if secretPath, ok := options[0].(string); ok {
					vaultOpts.SecretPath = secretPath
				}
			}
			WithVaultOptions(address, password, vaultOpts)(loader)

		case FileManager:
			WithFileOptions(address, nil)(loader)

		case EnvManager:
			WithEnvOptions(nil)(loader)

		case ConsulManager:
			consulOpts := &ConsulOptions{
				Username: username,
				Password: password,
			}
			WithConsulOptions(address, consulOpts)(loader)

		case EtcdManager:
			etcdOpts := &EtcdOptions{
				Username:  username,
				Password:  password,
				Endpoints: []string{address},
			}
			WithEtcdOptions([]string{address}, etcdOpts)(loader)

		case ZookeeperManager:
			zkOpts := &ZookeeperOptions{
				Username: username,
				Password: password,
				Servers:  []string{address},
			}
			WithZookeeperOptions([]string{address}, zkOpts)(loader)
		}
	}
}

// setConfigFields is a helper function that sets struct fields from a map using reflection
func setConfigFields(target interface{}, data map[string]interface{}, log Logger) error {
	// Get the type and value of the target struct
	v := reflect.ValueOf(target)

	// Check if target is a pointer to a struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	// Get the value to modify
	v = v.Elem()
	t := v.Type()

	// Create a map of lowercase env tags to their original values for case-insensitive lookup
	lowerCaseData := make(map[string]interface{})
	for key, val := range data {
		lowerCaseData[strings.ToLower(key)] = val
	}

	// Iterate through all fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the environment variable name from the tag
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// Convert the tag to lowercase for case-insensitive matching
		lowerEnvTag := strings.ToLower(envTag)

		// Check if the value exists in the data map (case-insensitive)
		value, exists := lowerCaseData[lowerEnvTag]
		if !exists {
			// Try the original tag as fallback
			value, exists = data[envTag]
			if !exists {
				continue
			}
		}

		// Set the field value
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		// Convert value to the appropriate type for the field
		switch fieldValue.Kind() {
		case reflect.String:
			strValue, ok := value.(string)
			if !ok {
				if value != nil {
					strValue = fmt.Sprintf("%v", value)
				} else {
					continue
				}
			}
			fieldValue.SetString(strValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var intValue int64
			switch v := value.(type) {
			case int:
				intValue = int64(v)
			case int64:
				intValue = v
			case float64:
				intValue = int64(v)
			case string:
				// Try to parse string to int
				// This is simplified and would need better parsing in a real implementation
				var i int
				if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
					intValue = int64(i)
				} else {
					continue
				}
			default:
				continue
			}
			fieldValue.SetInt(intValue)
		case reflect.Bool:
			boolValue, ok := value.(bool)
			if !ok {
				// Try to parse string to bool
				if strValue, ok := value.(string); ok {
					boolValue = strValue == "true" || strValue == "1" || strValue == "yes"
				} else {
					continue
				}
			}
			fieldValue.SetBool(boolValue)
		}
	}

	return nil
}

// Config is a type conversion helper
type Config config

// WithRetryOptions sets the retry options for all configuration managers
func WithRetryOptions(options *RetryOptions) ConfigOption {
	return func(loader *configLoader) {
		loader.retryOptions = options
		if loader.log != nil {
			extraMap := map[string]interface{}{
				"maxRetries":     options.MaxRetries,
				"initialBackoff": options.InitialBackoff.String(),
				"maxBackoff":     options.MaxBackoff.String(),
				"backoffFactor":  options.BackoffFactor,
			}
			loader.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Configuration, "Retry options configured", extraMap)
		}
	}
}
