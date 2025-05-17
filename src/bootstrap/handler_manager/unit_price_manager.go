package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	unitpriceusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/unit_price"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func UnitPriceInit(c contract.IContainer) *v1.UnitPriceHandler {
	use := unitpriceusecase.NewUnitPriceUsecase(c)
	handler := v1.NewUnitPriceHandler(use)

	return handler
}
