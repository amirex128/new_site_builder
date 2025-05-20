package sfconfigmanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

type vaultConfig struct {
	Address     string
	Token       string
	Timeout     time.Duration
	SecretPath  string
	AuthMethod  string
	Role        string
	SecretMount string
}

type vaultManager struct {
	client *api.Client
	config vaultConfig
	log    Logger
}

func NewVaultManager(config Config, log Logger) (configManager, error) {
	vaultCfg := vaultConfig{
		Address:     "http://127.0.0.1:8200",
		Timeout:     10 * time.Second,
		SecretPath:  "secret/data",
		SecretMount: "secret",
	}

	// Extract options from config
	if options, ok := config.Options["vault"].(map[string]interface{}); ok {
		if address, ok := options["address"].(string); ok && address != "" {
			vaultCfg.Address = address
		}
		if token, ok := options["token"].(string); ok && token != "" {
			vaultCfg.Token = token
		}
		if secretPath, ok := options["secretPath"].(string); ok && secretPath != "" {
			vaultCfg.SecretPath = secretPath
		}
		if secretMount, ok := options["secretMount"].(string); ok && secretMount != "" {
			vaultCfg.SecretMount = secretMount
		}
		if authMethod, ok := options["authMethod"].(string); ok && authMethod != "" {
			vaultCfg.AuthMethod = authMethod
		}
		if role, ok := options["role"].(string); ok && role != "" {
			vaultCfg.Role = role
		}
		if timeout, ok := options["timeout"].(time.Duration); ok && timeout > 0 {
			vaultCfg.Timeout = timeout
		}
	}

	// Create Vault client
	clientConfig := &api.Config{
		Address: vaultCfg.Address,
		Timeout: vaultCfg.Timeout,
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	// Set token if provided
	if vaultCfg.Token != "" {
		client.SetToken(vaultCfg.Token)
	}

	return &vaultManager{
		client: client,
		config: vaultCfg,
		log:    log,
	}, nil
}

func (vm *vaultManager) Load(target interface{}) (err error) {
	// Setup panic recovery
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = fmt.Errorf("vault manager panic: %s", x)
			case error:
				err = fmt.Errorf("vault manager panic: %w", x)
			default:
				err = fmt.Errorf("vault manager panic: %v", x)
			}

			if vm.log != nil {
				vm.log.ErrorWithCategory(
					Category.System.General,
					SubCategory.Status.Error,
					"Recovered from panic in Vault manager",
					map[string]interface{}{
						"error": err.Error(),
					},
				)
			}
		}
	}()

	// Get all secrets from the specified path
	secret, err := vm.client.Logical().Read(vm.config.SecretPath)
	if err != nil {
		return fmt.Errorf("failed to read from Vault: %w", err)
	}
	if secret == nil {
		return fmt.Errorf("no data found at path %s", vm.config.SecretPath)
	}

	// Extract data from the secret
	var data map[string]interface{}
	if vm.config.SecretMount == "secret" {
		// KV v2 path
		dataInterface, ok := secret.Data["data"]
		if !ok {
			return fmt.Errorf("no data found in secret")
		}
		data, ok = dataInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse secret data")
		}
	} else {
		// KV v1 path
		data = secret.Data
	}

	// Normalize keys to lowercase
	normalizedData := make(map[string]interface{})
	for k, v := range data {
		lowerKey := strings.ToLower(k)
		normalizedData[lowerKey] = v
	}

	vm.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Startup, "Loading configuration from Vault", map[string]interface{}{
		"address": vm.config.Address,
		"path":    vm.config.SecretPath,
		"count":   len(normalizedData),
	})

	// Set fields in the target struct
	return setConfigFields(target, normalizedData, vm.log)
}

func (vm *vaultManager) Get(key string) (interface{}, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	secret, err := vm.client.Logical().Read(vm.config.SecretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read from Vault: %w", err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no data found at path %s", vm.config.SecretPath)
	}

	var data map[string]interface{}
	if vm.config.SecretMount == "secret" {
		dataInterface, ok := secret.Data["data"]
		if !ok {
			return nil, fmt.Errorf("no data found in secret")
		}
		data, ok = dataInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse secret data")
		}
	} else {
		data = secret.Data
	}

	// Try first with lowercase key
	value, exists := data[lowerKey]
	if exists {
		return value, nil
	}

	// Try with original case as fallback
	value, exists = data[key]
	if !exists {
		// Check for any case-insensitive match
		for k, v := range data {
			if strings.ToLower(k) == lowerKey {
				return v, nil
			}
		}
		return nil, fmt.Errorf("key %s not found", key)
	}

	return value, nil
}

func (vm *vaultManager) GetString(key string) (string, error) {
	value, err := vm.Get(key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%v", value), nil
	}

	return strValue, nil
}

func (vm *vaultManager) GetInt(key string) (int, error) {
	value, err := vm.Get(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i, nil
		}
		return 0, fmt.Errorf("failed to convert %s to int", v)
	default:
		return 0, fmt.Errorf("cannot convert %v to int", value)
	}
}

func (vm *vaultManager) GetBool(key string) (bool, error) {
	value, err := vm.Get(key)
	if err != nil {
		return false, err
	}

	boolValue, ok := value.(bool)
	if ok {
		return boolValue, nil
	}

	if strValue, ok := value.(string); ok {
		return strValue == "true" || strValue == "1" || strValue == "yes", nil
	}

	return false, fmt.Errorf("cannot convert %v to bool", value)
}

func (vm *vaultManager) Close() error {
	// Nothing to close for Vault client
	return nil
}
