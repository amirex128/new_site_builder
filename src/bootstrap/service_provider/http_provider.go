package serviceprovider

import (
	"fmt"
	"git.snappfood.ir/backend/go/packages/sf-http-request/httpo"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
)

func HttpRequestProvider(logger sflogger.Logger) {
	err := httpo.RegisterConnection(
		httpo.WithConnectionDetails("badge", "https://api.example.com"),
		httpo.WithConnectionDetails("badge2", "https://api.example.com"),
		httpo.WithConnectionDetails("bid", "https://api.example.com"),
		httpo.WithLogger(logger),
	)
	if err != nil {
		logger.ErrorWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, fmt.Sprintf("Failed to register http connection: %v", err), nil)
	}
	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded Http", nil)

}
