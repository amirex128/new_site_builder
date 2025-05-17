package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	productcategoryusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product_category"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func ProductCategoryInit(c contract.IContainer) *v1.ProductCategoryHandler {
	use := productcategoryusecase.NewProductCategoryUsecase(c)
	handler := v1.NewProductCategoryHandler(use)

	return handler
}
