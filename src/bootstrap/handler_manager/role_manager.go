package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	roleusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/role"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func RoleInit(c contract.IContainer) *v1.RoleHandler {
	use := roleusecase.NewRoleUsecase(c)
	handler := v1.NewRoleHandler(use)

	return handler
}
