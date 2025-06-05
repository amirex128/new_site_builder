package consumerrouter

import (
	"context"
	"github.com/amirex128/new_site_builder/bootstrap"
	"github.com/amirex128/new_site_builder/internal/contract"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(ctx *context.Context, handlers *bootstrap.ConsumerHandlerManager, container contract.IContainer) {
	BindConsumers(ctx, handlers, container.GetLogger())

	// Block until SIGINT or SIGTERM is received
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	container.GetLogger().Infof("Shutdown signal received, exiting...")
}
