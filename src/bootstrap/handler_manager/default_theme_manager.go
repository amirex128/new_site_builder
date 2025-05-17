package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	defaultthemeusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/default_theme"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func DefaultThemeInit(c contract.IContainer) *v1.DefaultThemeHandler {
	use := defaultthemeusecase.NewDefaultThemeUsecase(c)
	handler := v1.NewDefaultThemeHandler(use)

	return handler
}
