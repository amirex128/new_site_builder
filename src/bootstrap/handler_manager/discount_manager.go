package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	discountusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/discount"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func DiscountInit(c contract.IContainer) *v1.DiscountHandler {
	use := discountusecase.NewDiscountUsecase(c)
	handler := v1.NewDiscountHandler(use)

	return handler
}
