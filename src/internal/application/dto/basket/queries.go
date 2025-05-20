package basket

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllBasketUserQuery for admin listing of all user baskets with pagination
type AdminGetAllBasketUserQuery struct {
	common.PaginationRequestDto
}

// GetAllBasketUserQuery for listing user baskets by site ID with pagination
type GetAllBasketUserQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" form:"siteId" validate:"required"`
}

// GetBasketQuery for retrieving a basket by site ID
type GetBasketQuery struct {
	SiteID *int64 `json:"siteId" form:"siteId" validate:"required"`
}
