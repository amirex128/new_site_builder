package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func LoggerProvider() sflogger.Logger {
	// Define logger options
	loggerOpts := []sflogger.Option{
		// Core logger configuration
		sflogger.WithLoggerType(sflogger.ZapLoggerType),
		sflogger.WithLevel(sflogger.InfoLevel),
		sflogger.WithAppName("SnappFood-Search"),
		sflogger.WithFormatter(sflogger.ColoredTextFormatter),
		sflogger.WithDevelopment(true),
		sflogger.WithStacktrace(true),

		// Configure log output destinations using helper functions
		// Console output for development visibility
		sflogger.WithConsoleSink(true), // true for colored output

		// File output with rotation
		sflogger.WithFileSink("./logs/app.log", 10, 30, 5, true),

		// Elasticsearch integration for centralized logging
		sflogger.WithMongoDBSink("localhost", 27017, "new_site_builder", "app_logs", "amirex128", "mI6G5jd3qNlJQinBOnA2z5SVEawLn4WV", 5, true),
	}

	// Initialize local logger
	logger := sflogger.New(loggerOpts...)

	// Initialize global logger with the same configuration
	sflogger.InitGlobalLogger(loggerOpts...)

	// Log with product_category-based logging
	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Application starting", map[string]interface{}{
		sflogger.ExtraKey.Metadata.AppName: "SnappFood-Search",
		sflogger.ExtraKey.Network.HostIP:   "127.0.0.1",
	})

	// Log using global logger
	sflogger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Global logger initialized", map[string]interface{}{
		sflogger.ExtraKey.Metadata.AppName: "SnappFood-Search",
	})
	return logger
}
