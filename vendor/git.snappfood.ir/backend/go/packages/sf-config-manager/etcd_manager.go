package sfconfigmanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdConfig struct {
	Endpoints []string
	Username  string
	Password  string
	Timeout   time.Duration
	BasePath  string
}

type etcdManager struct {
	client *clientv3.Client
	config etcdConfig
	log    Logger
}

func NewEtcdManager(config Config, log Logger) (configManager, error) {
	etcdCfg := etcdConfig{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "",
		Password:  "",
		Timeout:   10 * time.Second,
		BasePath:  "/config/",
	}

	// Extract options from config
	if options, ok := config.Options["etcd"].(map[string]interface{}); ok {
		if endpoints, ok := options["endpoints"].([]string); ok && len(endpoints) > 0 {
			etcdCfg.Endpoints = endpoints
		} else if endpoint, ok := options["endpoint"].(string); ok && endpoint != "" {
			etcdCfg.Endpoints = []string{endpoint}
		}
		if username, ok := options["username"].(string); ok {
			etcdCfg.Username = username
		}
		if password, ok := options["password"].(string); ok {
			etcdCfg.Password = password
		}
		if prefix, ok := options["prefix"].(string); ok && prefix != "" {
			etcdCfg.BasePath = prefix
			// Ensure prefix starts with "/"
			if !strings.HasPrefix(etcdCfg.BasePath, "/") {
				etcdCfg.BasePath = "/" + etcdCfg.BasePath
			}
		}
		if timeout, ok := options["timeout"].(time.Duration); ok && timeout > 0 {
			etcdCfg.Timeout = timeout
		}
	}

	// Create etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdCfg.Endpoints,
		Username:    etcdCfg.Username,
		Password:    etcdCfg.Password,
		DialTimeout: etcdCfg.Timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &etcdManager{
		client: client,
		config: etcdCfg,
		log:    log,
	}, nil
}

func (em *etcdManager) Load(target interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), em.config.Timeout)
	defer cancel()

	// Get all keys with the specified prefix
	resp, err := em.client.Get(ctx, em.config.BasePath, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to list keys from etcd: %w", err)
	}

	// Extract data from the key-value pairs
	data := make(map[string]interface{})
	for _, kv := range resp.Kvs {
		// Remove prefix from key
		key := strings.TrimPrefix(string(kv.Key), em.config.BasePath)
		// Skip empty keys
		if key == "" {
			continue
		}

		// Use dot notation for nested keys
		parts := strings.Split(key, "/")
		if len(parts) > 1 {
			// For nested keys, create nested maps
			createNestedMap(data, key, parts, string(kv.Value), em.log)
		} else {
			// Normalize key to lowercase
			lowerKey := strings.ToLower(key)
			data[lowerKey] = string(kv.Value)
		}
	}

	em.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Startup, "Loading configuration from etcd", map[string]interface{}{
		"endpoints": em.config.Endpoints,
		"basePath":  em.config.BasePath,
		"count":     len(resp.Kvs),
	})

	// Set fields in the target struct
	return setConfigFields(target, data, em.log)
}

func (em *etcdManager) Get(key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), em.config.Timeout)
	defer cancel()

	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Add prefix to key
	fullKey := em.config.BasePath + lowerKey

	// Get value from etcd
	resp, err := em.client.Get(ctx, fullKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get key from etcd: %w", err)
	}
	if len(resp.Kvs) == 0 {
		// Try with original case if lowercase key isn't found
		fullKey = em.config.BasePath + key
		resp, err = em.client.Get(ctx, fullKey)
		if err != nil || len(resp.Kvs) == 0 {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}

	return string(resp.Kvs[0].Value), nil
}

func (em *etcdManager) GetString(key string) (string, error) {
	value, err := em.Get(key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%v", value), nil
	}

	return strValue, nil
}

func (em *etcdManager) GetInt(key string) (int, error) {
	value, err := em.Get(key)
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

func (em *etcdManager) GetBool(key string) (bool, error) {
	value, err := em.Get(key)
	if err != nil {
		return false, err
	}

	strValue, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("value is not a string: %v", value)
	}

	return strValue == "true" || strValue == "1" || strValue == "yes", nil
}

func (em *etcdManager) Close() error {
	return em.client.Close()
}

func createNestedMap(keys map[string]interface{}, key string, parts []string, value interface{}, log Logger) {
	// Process nested keys
	current := keys
	for i, part := range parts[:len(parts)-1] {
		lowerPart := strings.ToLower(part)
		if _, ok := current[lowerPart]; !ok {
			current[lowerPart] = make(map[string]interface{})
		}

		nextMap, ok := current[lowerPart].(map[string]interface{})
		if !ok {
			// Cannot create nested structure, log a warning and return
			log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Cannot create nested structure for key %s at part %d", key, i), nil)
			return
		}

		current = nextMap
	}

	// Set the value at the leaf node
	lastPart := strings.ToLower(parts[len(parts)-1])
	current[lastPart] = value
}
