package consumerrouter

import (
	"context"
	"github.com/amirex128/new_site_builder/src/bootstrap"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func RunServer(ctx *context.Context, handlers *bootstrap.ConsumerHandlerManager, container contract.IContainer) {
	BindConsumers(ctx, handlers, container.GetLogger())
}
