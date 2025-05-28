package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/config"
)

func ElasticProvider(cfg *config.Config, logger sflogger.Logger) {
	//err := sfelastic.RegisterConnection(
	//	sfelastic.WithLogger(logger),
	//	sfelastic.WithGlobalOptions(func(client *elasticsearch.Client) {
	//	}),
	//	sfelastic.WithRetryOptions(&sfelastic.RetryOptions{
	//		MaxRetries:     5,                // Maximum number of retry attempts
	//		InitialBackoff: time.Second,      // Initial waiting time between retries
	//		MaxBackoff:     15 * time.Second, // Maximum waiting time between retries
	//		BackoffFactor:  1.5,              // Exponential backoff multiplier
	//	}),
	//	sfelastic.WithConnectionDetails("main", elasticsearch.Config{
	//		Addresses: []string{fmt.Sprintf("%s:%s", cfg.ElasticSearchHost, cfg.ElasticSearchPort)},
	//		Username:  cfg.ElasticSearchUsername,
	//		Password:  cfg.ElasticSearchPassword,
	//	}),
	//)
	//
	//if err != nil {
	//	logger.Errorf("Failed to register elasticsearch connection: %v", err)
	//}
}
