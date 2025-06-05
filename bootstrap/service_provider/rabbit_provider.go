package serviceprovider

import (
	"fmt"
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
	"github.com/amirex128/new_site_builder/config"
	"strconv"
	"time"
)

func RabbitProvider(cfg *config.Config, logger sflogger.Logger) {
	port, err := strconv.Atoi(cfg.RabbitmqPort)
	if err != nil {
		logger.Errorf("Could not convert RabbitMQ port to integer: %v", err)
	}
	err = sfrabbitmq.RegisterConnection(
		sfrabbitmq.WithLogger(logger),
		sfrabbitmq.WithGlobalOptions(func(c *sfrabbitmq.Config) {
			c.Heartbeat = 30 * time.Second
			c.Vhost = "/"
		}),
		sfrabbitmq.WithConnectionDetails(
			"main",
			cfg.RabbitmqHost,
			port,
			cfg.RabbitmqUsername,
			cfg.RabbitmqPassword,
		),
		// Configure outbox retry behavior
		sfrabbitmq.WithOutboxConfig(5*time.Second, 3),

		// Declare exchange
		sfrabbitmq.WithDeclareExchange("main", "notification", "direct"),
		// Declare queue
		sfrabbitmq.WithDeclareQueue("main", "sms"),
		sfrabbitmq.WithDeclareQueue("main", "email"),
		// Bind queue to exchange
		sfrabbitmq.WithBind("main", "sms", "send_sms", "notification"),
		sfrabbitmq.WithBind("main", "email", "send_email", "notification"),
	)

	if err != nil {
		logger.ErrorWithCategory(sflogger.Category.System.Startup, sflogger.SubCategory.Operation.Initialization, fmt.Sprintf("Failed to register rabbit connection: %v", err), nil)
	}
	logger.InfoWithCategory(sflogger.Category.System.General, sflogger.SubCategory.Operation.Startup, "Successfully loaded RabbitMQ", nil)

}
