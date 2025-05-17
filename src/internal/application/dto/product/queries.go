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
	SiteID *int64 `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// GetByIdProductQuery for retrieving a single product by ID
type GetByIdProductQuery struct {
	ID *int64 `json:"id" validate:"required" error:"required=محصول الزامی است"`
}

// GetSingleProductQuery for retrieving a single product by slug
type GetSingleProductQuery struct {
	Slug   *string `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SiteID *int64  `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// GetProductByCategoryQuery for retrieving products by product_category with pagination
type GetProductByCategoryQuery struct {
	common.PaginationRequestDto
	Slug   *string `json:"slug" validate:"required" error:"required=نامک الزامی است"`
	SiteID *int64  `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// GetByFiltersSortProductQuery for retrieving products with filtering and sorting
type GetByFiltersSortProductQuery struct {
	common.PaginationRequestDto
	SelectedFilters map[ProductFilterEnum][]string `json:"selectedFilters,omitempty" validate:"omitempty" error:""`
	SelectedSort    *ProductSortEnum               `json:"selectedSort,omitempty" validate:"omitempty" error:""`
	SiteID          *int64                         `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// CalculateProductsPriceQuery for calculating product prices
type CalculateProductsPriceQuery struct {
	CustomerID       *int64            `json:"customerId" validate:"required" error:"required=مشتری الزامی است"`
	SiteID           *int64            `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
	Code             *string           `json:"code,omitempty" validate:"omitempty" error:""`
	OrderBasketItems []OrderBasketItem `json:"orderBasketItems" validate:"required,dive" error:"required=آیتم‌های سبد خرید الزامی هستند"`
	IsOrderVerify    *bool             `json:"isOrderVerify" validate:"required" error:"required=وضعیت تایید سفارش الزامی است"`
}
