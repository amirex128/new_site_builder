package sfconfigmanager

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

type zookeeperConfig struct {
	Servers  []string
	Username string
	Password string
	Timeout  time.Duration
	BasePath string
}

type zookeeperManager struct {
	conn   *zk.Conn
	config zookeeperConfig
	log    Logger
}

func NewZookeeperManager(config Config, log Logger) (configManager, error) {
	zkCfg := zookeeperConfig{
		Servers:  []string{"127.0.0.1:2181"},
		Timeout:  10 * time.Second,
		BasePath: "/config",
	}

	// Extract options from config
	if options, ok := config.Options["zookeeper"].(map[string]interface{}); ok {
		if servers, ok := options["servers"].([]string); ok && len(servers) > 0 {
			zkCfg.Servers = servers
		} else if server, ok := options["server"].(string); ok && server != "" {
			zkCfg.Servers = []string{server}
		}
		if timeout, ok := options["sessionTimeout"].(time.Duration); ok && timeout > 0 {
			zkCfg.Timeout = timeout
		}
		if basePath, ok := options["basePath"].(string); ok && basePath != "" {
			zkCfg.BasePath = basePath
			// Ensure base path starts with "/"
			if !strings.HasPrefix(zkCfg.BasePath, "/") {
				zkCfg.BasePath = "/" + zkCfg.BasePath
			}
		}
	}

	// Connect to ZooKeeper
	conn, _, err := zk.Connect(zkCfg.Servers, zkCfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ZooKeeper: %w", err)
	}

	return &zookeeperManager{
		conn:   conn,
		config: zkCfg,
		log:    log,
	}, nil
}

func (zm *zookeeperManager) getChildren(path string) (map[string]interface{}, error) {
	// Get children of the current path
	children, _, err := zm.conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("failed to get children of %s: %w", path, err)
	}

	result := make(map[string]interface{})

	// Process each child
	for _, child := range children {
		childPath := path + "/" + child
		if strings.HasSuffix(childPath, "/") {
			// This is a directory, process its children recursively
			childData, err := zm.getChildren(childPath)
			if err != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				zm.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to process children of %s", childPath), extraMap)
				continue
			}
			result[child] = childData
		} else {
			// This is a leaf node, get its value
			data, _, err := zm.conn.Get(childPath)
			if err != nil {
				extraMap := map[string]interface{}{
					"error": err.Error(),
				}
				zm.log.WarnWithCategory(Category.System.General, SubCategory.Operation.Startup, fmt.Sprintf("Failed to get data for %s", childPath), extraMap)
				continue
			}
			result[child] = string(data)
		}
	}

	return result, nil
}

func (zm *zookeeperManager) Load(target interface{}) error {
	// Get all keys recursively from the base path
	data, err := zm.getChildren(zm.config.BasePath)
	if err != nil {
		return fmt.Errorf("failed to load configuration from ZooKeeper: %w", err)
	}
	if data == nil {
		return fmt.Errorf("base path %s not found", zm.config.BasePath)
	}

	zm.log.InfoWithCategory(Category.System.General, SubCategory.Operation.Startup, "Loading configuration from ZooKeeper", map[string]interface{}{
		"servers":  zm.config.Servers,
		"basePath": zm.config.BasePath,
		"count":    len(data),
	})

	// Set fields in the target struct
	return setConfigFields(target, data, zm.log)
}

func (zm *zookeeperManager) Get(key string) (interface{}, error) {
	// Convert key to lowercase for case-insensitive lookup
	lowerKey := strings.ToLower(key)

	// First try with the lowercase key
	fullPath := path.Join(zm.config.BasePath, lowerKey)
	data, _, err := zm.conn.Get(fullPath)
	if err != nil {
		if err == zk.ErrNoNode {
			// If not found with lowercase, try with original case
			fullPath = path.Join(zm.config.BasePath, key)
			data, _, err = zm.conn.Get(fullPath)
			if err != nil {
				if err == zk.ErrNoNode {
					return nil, fmt.Errorf("key not found: %s", key)
				}
				return nil, fmt.Errorf("failed to get key from ZooKeeper: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get key from ZooKeeper: %w", err)
		}
	}

	return string(data), nil
}

func (zm *zookeeperManager) GetString(key string) (string, error) {
	value, err := zm.Get(key)
	if err != nil {
		return "", err
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Sprintf("%v", value), nil
	}

	return strValue, nil
}

func (zm *zookeeperManager) GetInt(key string) (int, error) {
	value, err := zm.Get(key)
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

func (zm *zookeeperManager) GetBool(key string) (bool, error) {
	value, err := zm.Get(key)
	if err != nil {
		return false, err
	}

	strValue, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("value is not a string: %v", value)
	}

	return strValue == "true" || strValue == "1" || strValue == "yes", nil
}

func (zm *zookeeperManager) Close() error {
	zm.conn.Close()
	return nil
}
