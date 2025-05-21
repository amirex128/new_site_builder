package bootstrap

import "github.com/amirex128/new_site_builder/src/internal/contract"

type ConsumerHandlerManager struct {
}

func ConsumerHandlerBootstrap(c contract.IContainer) *ConsumerHandlerManager {
	return &ConsumerHandlerManager{}
}
