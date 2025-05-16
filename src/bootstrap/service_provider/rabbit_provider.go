package serviceprovider

import (
	sflogger "git.snappfood.ir/backend/go/packages/sf-logger"
	sfrabbitmq "git.snappfood.ir/backend/go/packages/sf-rabbitmq"
	"log"
	"time"
)

func RabbitProvider(logger sflogger.Logger) {
	err := sfrabbitmq.RegisterConnection(
		// Add the logger
		sfrabbitmq.WithLogger(logger),

		// Global options applied to all connections
		sfrabbitmq.WithGlobalOptions(func(c *sfrabbitmq.Config) {
			c.Heartbeat = 30 * time.Second
			c.Vhost = "/"
		}),

		// First connection with specific options
		sfrabbitmq.WithConnectionDetails(
			"order-connection",
			"localhost",
			5672,
			"guest",
			"guest",
		),

		sfrabbitmq.WithConnectionDetails(
			"ads-connection",
			"localhost",
			5672,
			"guest",
			"guest",
		),

		// Configure outbox retry behavior
		sfrabbitmq.WithOutboxConfig(5*time.Second, 3),

		// Declare exchange
		sfrabbitmq.WithDeclareExchange("order-connection", "order_exchange", "direct"),
		// Declare queue
		sfrabbitmq.WithDeclareQueue("order-connection", "order_queue"),
		// Bind queue to exchange
		sfrabbitmq.WithBind("order-connection", "order_queue", "order_routing_key", "order_exchange"),

		// Or Instead of the above three functions, you can only use the following function
		sfrabbitmq.WithDeclareExchangeAndQueue("ads-connection", "ads_exchange",
			"direct", "ads_queue", "ads_routing_key"),
	)

	if err != nil {
		log.Fatalf("Failed to register rabbit connection: %v", err)
	}

}
