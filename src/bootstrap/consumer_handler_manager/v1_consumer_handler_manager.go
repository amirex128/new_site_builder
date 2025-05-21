package consumerhandlermanager

import (
	"github.com/amirex128/new_site_builder/src/internal/api/handler/consumer"
	notificationusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/notification"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func NotificationInit(c contract.IContainer) *consumer.NotificationConsumerHandler {
	usecase := notificationusecase.NewNotificationUsecase(c)
	handler := consumer.NewNotificationConsumerHandler(usecase, c.GetLogger())
	return handler
}
