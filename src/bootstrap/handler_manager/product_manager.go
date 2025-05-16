package handlermanager

import (
	v1 "go-boilerplate/src/internal/api/handler/http/v1"
	productusecase "go-boilerplate/src/internal/app/usecase/product"
	"go-boilerplate/src/internal/contract"
)

func ProductInit(c contract.IContainer) *v1.ProductHandler {
	use := productusecase.NewProductUsecase(c)
	productList := v1.NewProductHandler(use)

	return productList
}
