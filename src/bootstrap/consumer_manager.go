package bootstrap

import (
	consumerhandlermanager "github.com/amirex128/new_site_builder/src/bootstrap/consumer_handler_manager"
	"github.com/amirex128/new_site_builder/src/internal/api/handler/consumer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

type ConsumerHandlerManager struct {
	NotificationHandler *consumer.NotificationConsumerHandler
}

func ConsumerHandlerBootstrap(c contract.IContainer) *ConsumerHandlerManager {

	return &ConsumerHandlerManager{
		NotificationHandler: consumerhandlermanager.NotificationInit(c),
	}
}
