package contract

// IConfig is the interface that provides configuration values
type IConfig interface {
	// GetString returns a string value for the given key
	GetString(key string) string

	// GetInt returns an int value for the given key
	GetInt(key string) int

	// GetBool returns a bool value for the given key
	GetBool(key string) bool

	// GetStringSlice returns a string slice for the given key
	GetStringSlice(key string) []string
}
