package serviceprovider

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"git.snappfood.ir/backend/go/packages/sf-mongo"
	"github.com/amirex128/new_site_builder/src/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func MongoProvider(cfg *config.Config, logger sflogger.Logger) {
	err := sfmongo.RegisterConnection(
		sfmongo.WithLogger(logger),
		sfmongo.WithConnectionDetails(
			"main",
			options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.MongodbUsername, cfg.MongodbPassword, cfg.MongodbHost, cfg.MongodbPort, cfg.MongodbDatabase)),
		),
		sfmongo.WithRetryOptions(&sfmongo.RetryOptions{
			MaxRetries:     5,
			InitialBackoff: time.Second,
			MaxBackoff:     15 * time.Second,
			BackoffFactor:  1.5,
		}),
	)
	if err != nil {
		logger.Errorf("Failed to register mongo MongoDB connection: %v", err)
	}
	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded Mongodb", nil)

}
