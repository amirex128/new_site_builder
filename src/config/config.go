package config

import (
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	AppLogLevel string `env:"APP_LOG_LEVEL"`
	AppPort     string `env:"APP_PORT"`

	MongodbDatabase string `env:"MONGODB_DATABASE"`
	MongodbHost     string `env:"MONGODB_HOST"`
	MongodbPassword string `env:"MONGODB_PASSWORD"`
	MongodbPort     string `env:"MONGODB_PORT"`
	MongodbUsername string `env:"MONGODB_USERNAME"`

	MysqlDatabase string `env:"MYSQL_DATABASE"`
	MysqlHost     string `env:"MYSQL_HOST"`
	MysqlPassword string `env:"MYSQL_PASSWORD"`
	MysqlPort     string `env:"MYSQL_PORT"`
	MysqlUsername string `env:"MYSQL_USERNAME"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisPort     string `env:"REDIS_PORT"`

	// Storage configuration
	StorageBucket    string `env:"STORAGE_BUCKET"`
	StorageRegion    string `env:"STORAGE_REGION"`
	StorageAccessKey string `env:"STORAGE_ACCESS_KEY"`
	StorageSecretKey string `env:"STORAGE_SECRET_KEY"`
}

// Storage returns storage configuration as a struct
func (c Config) Storage() struct {
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
} {
	return struct {
		Bucket    string
		Region    string
		AccessKey string
		SecretKey string
	}{
		Bucket:    c.StorageBucket,
		Region:    c.StorageRegion,
		AccessKey: c.StorageAccessKey,
		SecretKey: c.StorageSecretKey,
	}
}

func (c Config) GetString(key string) string {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				return field.String()
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				return field.String()
			}
			break
		}
	}

	return ""
}

func (c Config) GetInt(key string) int {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				if intValue, err := strconv.Atoi(field.String()); err == nil {
					return intValue
				}
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				if intValue, err := strconv.Atoi(field.String()); err == nil {
					return intValue
				}
			}
			break
		}
	}

	return 0
}

func (c Config) GetBool(key string) bool {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				if boolValue, err := strconv.ParseBool(field.String()); err == nil {
					return boolValue
				}
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				if boolValue, err := strconv.ParseBool(field.String()); err == nil {
					return boolValue
				}
			}
			break
		}
	}

	return false
}

func (c Config) GetStringSlice(key string) []string {
	// Convert the key to lowercase for case-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a case-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-insensitive)
		if strings.ToLower(envTag) == lowerKey {
			if field.IsValid() && field.Kind() == reflect.String {
				return []string{field.String()}
			}
			break
		}
	}

	// Fallback to exact tag name
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (case-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				return []string{field.String()}
			}
			break
		}
	}

	return nil
}
