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

	RabbitmqHost     string `env:"RABBITMQ_HOST"`
	RabbitmqPort     string `env:"RABBITMQ_PORT"`
	RabbitmqUsername string `env:"RABBITMQ_USERNAME"`
	RabbitmqPassword string `env:"RABBITMQ_PASSWORD"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisPort     string `env:"REDIS_PORT"`

	StorageS1Host      string `env:"STORAGE_S1_BUCKET"`
	StorageS1AccessKey string `env:"STORAGE_S1_ACCESS_KEY"`
	StorageS1SecretKey string `env:"STORAGE_S1_SECRET_KEY"`

	StorageS2Host      string `env:"STORAGE_S2_BUCKET"`
	StorageS2AccessKey string `env:"STORAGE_S2_ACCESS_KEY"`
	StorageS2SecretKey string `env:"STORAGE_S2_SECRET_KEY"`

	StorageS3Host      string `env:"STORAGE_S3_BUCKET"`
	StorageS3AccessKey string `env:"STORAGE_S3_ACCESS_KEY"`
	StorageS3SecretKey string `env:"STORAGE_S3_SECRET_KEY"`

	JwtSecretToken        string `env:"JWT_SECRET_TOKEN"`
	JwtIssuer             string `env:"JWT_ISSUER"`
	JwtAudience           string `env:"JWT_AUDIENCE"`
	ElasticSearchHost     string `env:"ELASTICSEARCH_HOST"`
	ElasticSearchPort     string `env:"ELASTICSEARCH_PORT"`
	ElasticSearchUsername string `env:"ELASTICSEARCH_USERNAME"`
	ElasticSearchPassword string `env:"ELASTICSEARCH_PASSWORD"`
}

func (c Config) GetString(key string) string {
	// Convert the key to lowercase for cache-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a cache-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (cache-insensitive)
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

		// Compare with requested key (cache-sensitive)
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
	// Convert the key to lowercase for cache-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a cache-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (cache-insensitive)
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

		// Compare with requested key (cache-sensitive)
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
	// Convert the key to lowercase for cache-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a cache-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (cache-insensitive)
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

		// Compare with requested key (cache-sensitive)
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
	// Convert the key to lowercase for cache-insensitive comparison
	lowerKey := strings.ToLower(key)

	// Get the value of the struct using reflection
	val := reflect.ValueOf(c)
	typ := val.Type()

	// Iterate through all fields to find a cache-insensitive match on env tag
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the env tag value
		envTag := fieldType.Tag.Get("env")

		// Compare with requested key (cache-insensitive)
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

		// Compare with requested key (cache-sensitive)
		if envTag == key {
			if field.IsValid() && field.Kind() == reflect.String {
				return []string{field.String()}
			}
			break
		}
	}

	return nil
}
