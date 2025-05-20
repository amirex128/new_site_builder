package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"git.snappfood.ir/backend/go/packages/sf-mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func MongoProvider(logger sflogger.Logger) {
	err := sfmongo.RegisterConnection(
		sfmongo.WithLogger(logger),
		sfmongo.WithConnectionDetails(
			"main-db",
			options.Client().ApplyURI("mongodb://localhost:27017"),
		),
		sfmongo.WithRetryOptions(&sfmongo.RetryOptions{
			MaxRetries:     5,
			InitialBackoff: time.Second,
			MaxBackoff:     15 * time.Second,
			BackoffFactor:  1.5,
		}),
	)
	if err != nil {
		logger.Errorf("Failed to register mongo database connection: %v", err)
	}

}
