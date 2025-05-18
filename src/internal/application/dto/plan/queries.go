package plan

import "github.com/amirex128/new_site_builder/src/internal/contract/common"

// GetAllPlanQuery represents a query to get all plans with pagination
type GetAllPlanQuery struct {
	common.PaginationRequestDto
}

// GetByIDPlanQuery represents a query to get a plan by ID
type GetByIDPlanQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0"`
}

// CalculatePlanPriceQuery represents a query to calculate plan price
type CalculatePlanPriceQuery struct {
	PlanID *int64 `json:"planId" validate:"required,gt=0"`
}
