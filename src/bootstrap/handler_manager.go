package bootstrap

import (
	handlermanager "go-boilerplate/src/bootstrap/handler_manager"
	"go-boilerplate/src/internal/api/handler/http/v1"
	"go-boilerplate/src/internal/contract"
)

func HandlerBootstrap(c contract.IContainer) *HandlerManager {

	return &HandlerManager{
		ProductHandlerV1: handlermanager.ProductInit(c),
		VendorHandlerV1:  handlermanager.VendorInit(c),
	}
}

type HandlerManager struct {
	ProductHandlerV1 *v1.ProductHandler
	VendorHandlerV1  *v1.VendorHandler
	//declare other handlers
}
