package sfconfigmanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

type consulConfig struct {
	Address    string
	Datacenter string
	Token      string
	Timeout    time.Duration
	Prefix     string
}

type consulManager struct {
	client *api.Client
	config consulConfig
	log    Logger
}

func NewConsulManager(config Config, log Logger) (configManager, error) {
	consulCfg := consulConfig{
		Address:    "127.0.0.1:8500",
		Datacenter: "",
		Token:      "",
		Timeout:    10 * time.Second,
		Prefix:     "config/",
	}

	// Extract options from config
	if options, ok := config.Options["consul"].(map[string]interface{}); ok {
		if address, ok := options["address"].(string); ok && address != "" {
			consulCfg.Address = address
		}
		if datacenter, ok := options["datacenter"].(string); ok && datacenter != "" {
			consulCfg.Datacenter = datacenter
		}
		if token, ok := options["token"].(string); ok && token != "" {
			consulCfg.Token = token
		}
		if prefix, ok := options["prefix"].(string); ok && prefix != "" {
			consulCfg.Prefix = prefix
		}
		if timeout, ok := options["timeout"].(time.Duration); ok && timeout > 0 {
			consulCfg.Timeout = timeout
		}
	}

	// Create Consul client
	clientConfig := &api.Config{
		Address:    consulCfg.Address,
		Datacenter: consulCfg.Datacenter,
		Token:      consulCfg.Token,
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &consulManager{
		client: client,
		config: consulCfg,
		log:    log,
	}, nil
}

func (cm *consulManager) Load(target interface{}) error {
	// Get all keys from Consul KV store with the specified prefix
	kv := cm.client.KV()
	pairs, _, err := kv.List(cm.config.Prefix, nil)
	if err != nil {
		return fmt.Errorf("failed to list keys from Consul: %w", err)
	}

	// Extract data from the key-value pairs
	data := make(map[string]interface{})
	for _, pair := range pairs {
		// Remove prefix from key
		key := strings.TrimPrefix(pair.Key, cm.config.Prefix)
		// Skip empty keys or directory entries
		if key == "" || pair.Value == nil {
			continue
		}

		// Use dot notation for nested keys
		parts := strings.Split(key, "/")
		if len(parts) > 1 {
			// For nested keys, create nested maps
			current := data
			for i, part := range parts[:len(parts)-1] {
				if part == "" {
					continue
				}
				// Normalize part to lowercase
				lowerPart := strings.ToLower(part)
				if _, exists := current[lowerPart]; !exists {
					current[lowerPart] = make(map[string]interface{})
				}
				if nestedMap, ok := current[lowerPart].(map[string]interface{}); ok {
					current = nestedMap
				} else {
					cm.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Cannot create nested structure for key %s at part %d", key, i), nil)
					break
				}
			}
			lastPart := parts[len(parts)-1]
			if lastPart != "" {
				// Normalize the last part to lowercase
				lowerLastPart := strings.ToLower(lastPart)
				current[lowerLastPart] = string(pair.Value)
			}
		} else {
			// Normalize the key to lowercase
			lowerKey := strings.ToLower(key)
			data[lowerKey] = string(pair.Value)
		}
	}

	cm.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Startup, "Loading configuration from Consul", map[string]interface{}{
		"address": cm.config.Address,
		"prefix":  cm.config.Prefix,
		"count":   len(pairs),
	})

	// Set fields in the target struct
	return setConfigFields(target, data, cm.log)
}

func (cm *consulManager) Get(key string) (interface{}, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Add prefix to key
	fullKey := cm.config.Prefix + lowerKey

	// Get value from Consul KV store
	kv := cm.client.KV()
	pair, _, err := kv.Get(fullKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get key from Consul: %w", err)
	}
	if pair == nil {
		// Try with original case if lowercase key isn't found
		fullKey = cm.config.Prefix + key
		pair, _, err = kv.Get(fullKey, nil)
		if err != nil || pair == nil {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}

	return string(pair.Value), nil
}

func (cm *consulManager) GetString(key string) (string, error) {
	value, err := cm.Get(key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%v", value), nil
	}

	return strValue, nil
}

func (cm *consulManager) GetInt(key string) (int, error) {
	value, err := cm.Get(key)
	if err != nil {
		return 0, err
	}

	strValue, ok := value.(string)
	if !ok {
		return 0, fmt.Errorf("value is not a string: %v", value)
	}

	var i int
	if _, err := fmt.Sscanf(strValue, "%d", &i); err != nil {
		return 0, fmt.Errorf("failed to convert %s to int: %w", strValue, err)
	}

	return i, nil
}

func (cm *consulManager) GetBool(key string) (bool, error) {
	value, err := cm.Get(key)
	if err != nil {
		return false, err
	}

	strValue, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("value is not a string: %v", value)
	}

	return strValue == "true" || strValue == "1" || strValue == "yes", nil
}

func (cm *consulManager) Close() error {
	// Nothing to close for Consul client
	return nil
}
