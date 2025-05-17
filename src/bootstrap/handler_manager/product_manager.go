package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	productusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func ProductInit(c contract.IContainer) *v1.ProductHandler {
	use := productusecase.NewProductUsecase(c)
	handler := v1.NewProductHandler(use)

	return handler
}
