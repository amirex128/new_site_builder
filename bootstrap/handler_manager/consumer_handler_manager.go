package handlermanager

import (
	"github.com/amirex128/new_site_builder/internal/api/handler/consumer"
	"github.com/amirex128/new_site_builder/internal/application/usecase/notification"
	"github.com/amirex128/new_site_builder/internal/contract"
)

func NotificationInit(c contract.IContainer) *consumer.NotificationConsumerHandler {
	usecase := notificationusecase.NewNotificationUsecase(c)
	handler := consumer.NewNotificationConsumerHandler(usecase, c.GetLogger())
	return handler
}
