package plan

import "github.com/amirex128/new_site_builder/src/internal/contract/common"

// GetAllPlanQuery represents a query to get all plans with pagination
type GetAllPlanQuery struct {
	common.PaginationRequestDto
}

// GetByIDPlanQuery represents a query to get a plan by ID
type GetByIDPlanQuery struct {
	ID *int64 `json:"id" form:"id" validate:"required,gt=0"`
}

// CalculatePlanPriceQuery represents a query to calculate plan price
type CalculatePlanPriceQuery struct {
	PlanID *int64 `json:"planId" form:"planId" validate:"required,gt=0"`
}

// PlanSiteID represents the site ID for plan queries
type PlanSiteID struct {
	SiteID *int64 `json:"siteId" nameFa:"شناسه سایت" form:"siteId" validate:"required"`
}
