package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	siteusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/site"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func SiteInit(c contract.IContainer) *v1.SiteHandler {
	use := siteusecase.NewSiteUsecase(c)
	handler := v1.NewSiteHandler(use)

	return handler
}
