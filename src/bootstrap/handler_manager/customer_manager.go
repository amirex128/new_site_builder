package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	customerusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/customer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func CustomerInit(c contract.IContainer) *v1.CustomerHandler {
	use := customerusecase.NewCustomerUsecase(c)
	handler := v1.NewCustomerHandler(use)

	return handler
}
