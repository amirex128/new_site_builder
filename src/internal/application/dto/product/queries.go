package product

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllProductQuery for admin listing of all products with pagination
type AdminGetAllProductQuery struct {
	common.PaginationRequestDto
}

// GetAllProductQuery for listing products by site ID with pagination
type GetAllProductQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" form:"siteId" validate:"required"`
}

// GetByIdProductQuery for retrieving a single article by ID
type GetByIdProductQuery struct {
	ID *int64 `json:"id" form:"id" validate:"required"`
}

// GetSingleProductQuery for retrieving a single article by slug
type GetSingleProductQuery struct {
	Slug   *string `json:"slug" form:"slug" validate:"required"`
	SiteID *int64  `json:"siteId" form:"siteId" validate:"required"`
}

// GetProductByCategoryQuery for retrieving products by product_category with pagination
type GetProductByCategoryQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" form:"slug" validate:"required"`
	SiteID *int64  `json:"siteId" form:"siteId" validate:"required"`
}

// GetByFiltersSortProductQuery for retrieving products with filtering and sorting
type GetByFiltersSortProductQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[ProductFilterEnum][]string `json:"selectedFilters,omitempty" form:"selectedFilters" validate:"omitempty"`
	SelectedSort    *ProductSortEnum               `json:"selectedSort,omitempty" form:"selectedSort" validate:"omitempty"`
	SiteID          *int64                         `json:"siteId" form:"siteId" validate:"required"`
}

// CalculateProductsPriceQuery for calculating article prices
type CalculateProductsPriceQuery struct {
	CustomerID       *int64            `json:"customerId" validate:"required"`
	SiteID           *int64            `json:"siteId" validate:"required"`
	Code             *string           `json:"code,omitempty" validate:"omitempty"`
	OrderBasketItems []OrderBasketItem `json:"orderBasketItems" validate:"required,dive"`
	IsOrderVerify    *bool             `json:"isOrderVerify" validate:"required"`
}
