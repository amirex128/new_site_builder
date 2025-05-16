package sfconfigmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type fileConfig struct {
	Path string
	Type string // json, yaml
}

type fileManager struct {
	config fileConfig
	data   map[string]interface{}
	log    Logger
}

func NewFileManager(config Config, log Logger) (configManager, error) {
	fileCfg := fileConfig{
		Path: "",
		Type: "json",
	}

	// Extract options from config
	if options, ok := config.Options["file"].(map[string]interface{}); ok {
		if path, ok := options["path"].(string); ok && path != "" {
			fileCfg.Path = path
		}
		if fileType, ok := options["type"].(string); ok && fileType != "" {
			fileCfg.Type = strings.ToLower(fileType)
		}
	}

	if fileCfg.Path == "" {
		return nil, fmt.Errorf("file path is required")
	}

	// Check if file type is supported
	if fileCfg.Type != "json" && fileCfg.Type != "yaml" && fileCfg.Type != "yml" {
		return nil, fmt.Errorf("unsupported file type: %s", fileCfg.Type)
	}

	// Create file manager
	fm := &fileManager{
		config: fileCfg,
		data:   make(map[string]interface{}),
		log:    log,
	}

	// Load configuration data
	if err := fm.loadFile(); err != nil {
		return nil, err
	}

	return fm, nil
}

// Helper function to convert all map keys to lowercase recursively
func normalizeMapKeys(data map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})

	for k, v := range data {
		// Convert current key to lowercase
		lowerKey := strings.ToLower(k)

		// Check if value is a nested map and recursively normalize it
		if nestedMap, ok := v.(map[string]interface{}); ok {
			normalized[lowerKey] = normalizeMapKeys(nestedMap)
		} else if nestedMap, ok := v.(map[interface{}]interface{}); ok {
			// Convert map[interface{}]interface{} to map[string]interface{}
			stringMap := make(map[string]interface{})
			for mk, mv := range nestedMap {
				if strKey, ok := mk.(string); ok {
					stringMap[strings.ToLower(strKey)] = mv
				}
			}
			normalized[lowerKey] = normalizeMapKeys(stringMap)
		} else {
			normalized[lowerKey] = v
		}
	}

	return normalized
}

func (fm *fileManager) loadFile() error {
	// Check if file exists
	if _, err := os.Stat(fm.config.Path); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", fm.config.Path)
	}

	// Read file content
	content, err := os.ReadFile(fm.config.Path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse file content based on file type
	var data map[string]interface{}
	switch fm.config.Type {
	case "json":
		if err := json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
	}

	// Normalize all keys to lowercase
	fm.data = normalizeMapKeys(data)
	return nil
}

// Load configuration data from the file
func (fm *fileManager) Load(target interface{}) error {
	fm.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Startup, "Loading configuration from file", map[string]interface{}{
		"path": fm.config.Path,
		"type": fm.config.Type,
	})

	// Set fields in the target struct
	return setConfigFields(target, fm.data, fm.log)
}

func (fm *fileManager) Get(key string) (interface{}, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// Support nested keys with dot notation (e.g. "database.host")
	keys := strings.Split(lowerKey, ".")
	var value interface{} = fm.data

	for _, k := range keys {
		// Try to navigate to the nested key
		if mapValue, ok := value.(map[string]interface{}); ok {
			if v, exists := mapValue[k]; exists {
				value = v
			} else {
				return nil, fmt.Errorf("key not found: %s", key)
			}
		} else if mapValue, ok := value.(map[interface{}]interface{}); ok {
			// Support for YAML with non-string keys
			found := false
			for mk, mv := range mapValue {
				if strKey, ok := mk.(string); ok && strings.ToLower(strKey) == k {
					value = mv
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("key not found: %s", key)
			}
		} else {
			return nil, fmt.Errorf("cannot access nested key: %s", key)
		}
	}

	return value, nil
}

func (fm *fileManager) GetString(key string) (string, error) {
	value, err := fm.Get(key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%v", value), nil
	}

	return strValue, nil
}

func (fm *fileManager) GetInt(key string) (int, error) {
	value, err := fm.Get(key)
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

func (fm *fileManager) GetBool(key string) (bool, error) {
	value, err := fm.Get(key)
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

func (fm *fileManager) Close() error {
	// Nothing to close for file manager
	return nil
}
