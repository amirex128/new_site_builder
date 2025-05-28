package serviceprovider

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfredis "git.snappfood.ir/backend/go/packages/sf-redis"
	"github.com/amirex128/new_site_builder/src/config"
	"time"
)

func RedisProvider(cfg *config.Config, logger sflogger.Logger) {
	err := sfredis.RegisterConnection(
		sfredis.WithRetryOptions(&sfredis.RetryOptions{
			MaxRetries:     5,
			InitialBackoff: time.Second,
			MaxBackoff:     15 * time.Second,
			BackoffFactor:  1.5,
		}),
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
		logger.ErrorWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, fmt.Sprintf("Failed to register redis connection: %v", err), nil)
	}

	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded Redis", nil)

}
