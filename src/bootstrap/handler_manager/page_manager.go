package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	pageusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/page"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func PageInit(c contract.IContainer) *v1.PageHandler {
	use := pageusecase.NewPageUsecase(c)
	handler := v1.NewPageHandler(use)

	return handler
}
