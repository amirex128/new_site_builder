package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	addressusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/address"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func AddressInit(c contract.IContainer) *v1.AddressHandler {
	use := addressusecase.NewAddressUsecase(c)
	handler := v1.NewAddressHandler(use)

	return handler
}
