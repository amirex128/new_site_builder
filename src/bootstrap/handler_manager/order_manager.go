package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	orderusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/order"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func OrderInit(c contract.IContainer) *v1.OrderHandler {
	use := orderusecase.NewOrderUsecase(c)
	handler := v1.NewOrderHandler(use)

	return handler
}
