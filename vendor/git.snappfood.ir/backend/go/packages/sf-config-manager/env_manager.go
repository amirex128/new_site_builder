package sfconfigmanager

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type envConfig struct {
	Prefix string
}

type envManager struct {
	config envConfig
	log    Logger
}

func NewEnvManager(config Config, log Logger) (configManager, error) {
	envCfg := envConfig{
		Prefix: "",
	}

	// Extract options from config
	if options, ok := config.Options["env"].(map[string]interface{}); ok {
		if prefix, ok := options["prefix"].(string); ok {
			envCfg.Prefix = prefix
		}
	}

	return &envManager{
		config: envCfg,
		log:    log,
	}, nil
}

func (em *envManager) Load(target interface{}) error {
	// Get all environment variables
	envVars := make(map[string]interface{})

	// For this implementation, we'll use environment variables directly
	for _, key := range os.Environ() {
		parts := strings.SplitN(key, "=", 2)
		if len(parts) == 2 {
			envVars[strings.ToLower(parts[0])] = parts[1]
		}
	}

	em.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Startup, "Loading configuration from environment variables", map[string]interface{}{
		"prefix": em.config.Prefix,
		"count":  len(envVars),
	})

	// Set fields in the target struct
	return setConfigFields(target, envVars, em.log)
}

func (em *envManager) Get(key string) (interface{}, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Try with lowercase key first
	if os.Getenv(lowerKey) != "" {
		return os.Getenv(lowerKey), nil
	}

	// Try with original key as fallback
	if os.Getenv(key) != "" {
		return os.Getenv(key), nil
	}

	return nil, fmt.Errorf("environment variable %s not found", key)
}

func (em *envManager) GetString(key string) (string, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Try with lowercase key first
	if os.Getenv(lowerKey) != "" {
		return os.Getenv(lowerKey), nil
	}

	// Try with original key as fallback
	if os.Getenv(key) != "" {
		return os.Getenv(key), nil
	}

	return "", fmt.Errorf("environment variable %s not found", key)
}

func (em *envManager) GetInt(key string) (int, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Try with lowercase key first
	if os.Getenv(lowerKey) != "" {
		return strconv.Atoi(os.Getenv(lowerKey))
	}

	// Try with original key as fallback
	if os.Getenv(key) != "" {
		return strconv.Atoi(os.Getenv(key))
	}

	return 0, fmt.Errorf("environment variable %s not found", key)
}

func (em *envManager) GetBool(key string) (bool, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Try with lowercase key first
	if os.Getenv(lowerKey) != "" {
		return strconv.ParseBool(os.Getenv(lowerKey))
	}

	// Try with original key as fallback
	if os.Getenv(key) != "" {
		return strconv.ParseBool(os.Getenv(key))
	}

	return false, fmt.Errorf("environment variable %s not found", key)
}

func (em *envManager) Close() error {
	// Nothing to close for environment variables
	return nil
}
