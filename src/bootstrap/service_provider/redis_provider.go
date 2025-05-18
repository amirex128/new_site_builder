package serviceprovider

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
	"github.com/amirex128/new_site_builder/src/config"
	"log"
)

func RedisProvider(cfg *config.Config, logger sflogger.Logger) {
	err := sfredis.RegisterConnection(
		sfredis.WithGlobalOptions(func(options *sfredis.Options) {
			options.PoolSize = 20
			options.MinIdleConns = 5
			options.MaxRetries = 3
		}),
		// Connection with custom options
		sfredis.WithConnectionDetails("cache", fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort), cfg.RedisPassword, 0,
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
