package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	headerfooterusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/header_footer"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func HeaderFooterInit(c contract.IContainer) *v1.HeaderFooterHandler {
	use := headerfooterusecase.NewHeaderFooterUsecase(c)
	handler := v1.NewHeaderFooterHandler(use)

	return handler
}
