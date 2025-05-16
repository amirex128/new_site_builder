package serviceprovider

import (
	elastic "git.snappfood.ir/backend/go/packages/sf-elasticsearch-client/v8"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func ElasticProvider(logger sflogger.Logger) {
	// Register the Elasticsearch connection with the logger
	err := elastic.RegisterConnection(
		elastic.WithLogger(logger),
		elastic.WithGlobalOptions(func(client *elasticsearch.Client) {
			// set global option
		}),
		elastic.WithConnectionDetails("main", elasticsearch.Config{
			Addresses: []string{"http://localhost:9200"},
			Username:  "elastic",
			Password:  "changeme",
		}),
	)

	if err != nil {
		log.Fatalf("Failed to register elastic connection: %v", err)
	}
}
