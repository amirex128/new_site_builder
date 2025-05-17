package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	userusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/user"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func UserInit(c contract.IContainer) *v1.UserHandler {
	use := userusecase.NewUserUsecase(c)
	handler := v1.NewUserHandler(use)

	return handler
}
