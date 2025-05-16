package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
	"log"
)

func RedisProvider(logger sflogger.Logger) {
	err := sfredis.RegisterConnection(
		sfredis.WithGlobalOptions(func(options *sfredis.Options) {
			options.PoolSize = 20
			options.MinIdleConns = 5
			options.MaxRetries = 3
		}),

		// Connection using only global options
		sfredis.WithConnectionDetails("session", "127.0.0.1:6399", "", 0),

		// Connection with custom options
		sfredis.WithConnectionDetails("cache", "127.0.0.1:6399", "", 0,
			func(options *sfredis.Options) {
				options.PoolSize = 30
				options.MinIdleConns = 10
				options.MaxRetries = 5
			},
		),
	)
	if err != nil {
		log.Fatalf("Failed to register redis connection: %v", err)
	}
}
