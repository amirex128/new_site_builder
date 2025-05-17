package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	basketusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/basket"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func BasketInit(c contract.IContainer) *v1.BasketHandler {
	use := basketusecase.NewBasketUsecase(c)
	handler := v1.NewBasketHandler(use)

	return handler
}
