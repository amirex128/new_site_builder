package serviceprovider

import (
	"fmt"
	elastic "git.snappfood.ir/backend/go/packages/sf-elasticsearch-client/v8"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/amirex128/new_site_builder/src/config"
	"github.com/elastic/go-elasticsearch/v8"
)

func ElasticProvider(cfg *config.Config, logger sflogger.Logger) {
	err := elastic.RegisterConnection(
		elastic.WithLogger(logger),
		elastic.WithGlobalOptions(func(client *elasticsearch.Client) {
		}),
		elastic.WithConnectionDetails("main", elasticsearch.Config{
			Addresses: []string{fmt.Sprintf("%s:%s", cfg.ElasticSearchHost, cfg.ElasticSearchPort)},
			Username:  cfg.ElasticSearchUsername,
			Password:  cfg.ElasticSearchPassword,
		}),
	)

	if err != nil {
		logger.Errorf("Failed to register elasticsearch connection: %v", err)
	}
}
