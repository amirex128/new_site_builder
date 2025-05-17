package bootstrap

import (
	handlermanager "github.com/amirex128/new_site_builder/src/bootstrap/handler_manager"
	"github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	"github.com/amirex128/new_site_builder/src/internal/contract"
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
