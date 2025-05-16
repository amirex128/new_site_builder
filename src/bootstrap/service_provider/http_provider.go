package serviceprovider

import (
	"git.snappfood.ir/backend/go/packages/sf-http-request/httpo"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	"log"
)

func HttpProvider(logger sflogger.Logger) {
	err := httpo.RegisterConnection(
		httpo.WithConnectionDetails("badge", "https://api.example.com"),
		httpo.WithConnectionDetails("badge2", "https://api.example.com"),
		httpo.WithConnectionDetails("bid", "https://api.example.com"),
		httpo.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to register connection: %v", err)
	}
}
