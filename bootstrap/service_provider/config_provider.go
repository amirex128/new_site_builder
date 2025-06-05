package serviceprovider

import (
	"fmt"
	"github.com/amirex128/new_site_builder/config"
	"os"
	"time"

	sfconfigmanager "git.snappfood.ir/backend/go/packages/sf-config-manager"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func ConfigProvider(logger sflogger.Logger) *config.Config {
	cfg := &config.Config{}
	result, err := sfconfigmanager.RegisterConnection(
		sfconfigmanager.WithConfig(cfg),
		sfconfigmanager.WithLogger(logger),
		sfconfigmanager.WithRetryOptions(&sfconfigmanager.RetryOptions{
			MaxRetries:     5,                // Maximum number of retry attempts
			InitialBackoff: time.Second,      // Initial waiting time between retries
			MaxBackoff:     15 * time.Second, // Maximum waiting time between retries
			BackoffFactor:  1.5,              // Exponential backoff multiplier
		}),
		sfconfigmanager.WithFileOptions(
			getConfigPath(os.Getenv("APP_ENV")),
			&sfconfigmanager.FileOptions{
				Type: "json",
			},
		),

		sfconfigmanager.WithEnvOptions(nil),
	)

	if err != nil {
		logger.ErrorWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, fmt.Sprintf("Failed to load configuration from any source : %v", err), nil)
		return nil
	}

	cfg, ok := result.(*config.Config)
	if !ok {
		logger.ErrorWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Failed to cast configuration result to config.Config", nil)
		return nil
	}

	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded configuration", nil)
	return cfg
}

func getConfigPath(env string) string {
	if env == "local" {
		return "config/files/config-local.json"
	} else if env == "stage" {
		return "config/files/config-stage.json"
	} else if env == "production" {
		return "src/config/files/config-production.json"
	}
	return "config/files/config-local.json"
}
