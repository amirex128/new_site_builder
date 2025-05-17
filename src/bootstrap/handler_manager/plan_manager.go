package handlermanager

import (
	v1 "github.com/amirex128/new_site_builder/src/internal/api/handler/http/v1"
	planusecase "github.com/amirex128/new_site_builder/src/internal/application/usecase/plan"
	"github.com/amirex128/new_site_builder/src/internal/contract"
)

func PlanInit(c contract.IContainer) *v1.PlanHandler {
	use := planusecase.NewPlanUsecase(c)
	handler := v1.NewPlanHandler(use)

	return handler
}
