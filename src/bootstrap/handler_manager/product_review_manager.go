package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	productreviewusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/product_review"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func ProductReviewInit(c contract.IContainer) *v1.ProductReviewHandler {
	use := productreviewusecase.NewProductReviewUsecase(c)
	handler := v1.NewProductReviewHandler(use)

	return handler
}
