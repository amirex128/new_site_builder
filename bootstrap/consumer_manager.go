package bootstrap

import (
	"github.com/amirex128/new_site_builder/bootstrap/handler_manager"
	"github.com/amirex128/new_site_builder/internal/api/handler/consumer"
	"github.com/amirex128/new_site_builder/internal/contract"
)

type ConsumerHandlerManager struct {
	NotificationHandler *consumer.NotificationConsumerHandler
}

func ConsumerHandlerBootstrap(c contract.IContainer) *ConsumerHandlerManager {

	return &ConsumerHandlerManager{
		NotificationHandler: handlermanager.NotificationInit(c),
	}
}
