package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	websiteusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/website"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func WebsiteInit(c contract.IContainer) *v1.WebsiteHandler {
	use := websiteusecase.NewWebsiteUsecase(c)
	handler := v1.NewWebsiteHandler(use)

	return handler
}
