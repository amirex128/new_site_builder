package serviceprovider

import (
	"github.com/amirex128/new_site_builder/src/config"
	"os"

	sfconfigmanager "git.snappfood.ir/backend/go/packages/sf-config-manager"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

// ConfigProvider provides configuration from various sources with fallback
func ConfigProvider(logger sflogger.Logger) *config.Config {
	// Create a struct to hold the configuration
	cfg := &config.Config{}

	// Register connection with fallback sources in order:
	// 1. Vault
	// 2. File
	// 3. Environment Variables
	result, err := sfconfigmanager.RegisterConnection(
		sfconfigmanager.WithConfig(cfg),
		sfconfigmanager.WithLogger(logger),

		// First try Vault
		sfconfigmanager.WithVaultOptions(
			"http://localhost:8200",
			"mI6G5jd3qNlJQinBOnA2z5SVEawLn4WV",
			&sfconfigmanager.VaultOptions{
				SecretPath:  getVaultSecretPath(),
				SecretMount: "secret",
			},
		),

		// Then try file configuration
		sfconfigmanager.WithFileOptions(
			getConfigPath(os.Getenv("APP_ENV")),
			&sfconfigmanager.FileOptions{
				Type: "yml",
			},
		),

		// Finally try environment variables
		sfconfigmanager.WithEnvOptions(nil),
	)

	if err != nil {
		// All configuration sources failed
		extraMap := map[string]interface{}{
			sflogger.ExtraKey.Error.ErrorMessage: err.Error(),
		}
		logger.ErrorWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Failed to load configuration from any source", extraMap)
		return nil
	}

	// Cast the result back to config.Config
	config, ok := result.(*config.Config)
	if !ok {
		logger.ErrorWithCategory(sflogger.Category.Database.Redis, sflogger.SubCategory.Operation.Startup, "Failed to cast configuration result to config.Config", nil)
		return nil
	}

	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded configuration", nil)
	return config
}

// getVaultSecretPath returns the secret path for Vault
func getVaultSecretPath() string {
	path := os.Getenv("VAULT_SECRET_PATH")
	if path == "" {
		return "kv/data/new_site_builder" // Default path
	}
	return path
}

// getConfigPath returns the path to the config file based on environment
func getConfigPath(env string) string {
	if env == "docker" {
		return "/app/config/config-docker.yml"
	} else if env == "production" {
		return "/config/config-production.yml"
	}

	return "src/config/config-development.yml"
}
