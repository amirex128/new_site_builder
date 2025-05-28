package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func LoggerProvider() sflogger.Logger {
	logger := sflogger.RegisterLogger(
		sflogger.WithAppName("my-appaaa"),
		sflogger.WithLevel(sflogger.InfoLevel),
		sflogger.WithFormatter(sflogger.ColoredTextFormatter),
		sflogger.WithStacktrace(true),
		sflogger.WithMongoDBSink(
			"localhost",
			27017,
			"new_site_builder",
			"app_logs",
			"amirex128",
			"mI6G5jd3qNlJQinBOnA2z5SVEawLn4WV",
			5,
			false,
		),
	)

	logger.InfoWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, "Logger initialized successfully", nil)

	return logger
}
