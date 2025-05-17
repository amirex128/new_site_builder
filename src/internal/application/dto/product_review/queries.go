package product_review

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllProductReviewQuery for admin listing of all product reviews with pagination
type AdminGetAllProductReviewQuery struct {
	common.PaginationRequestDto
}

// GetAllProductReviewQuery for listing product reviews by site ID with pagination
type GetAllProductReviewQuery struct {
	common.PaginationRequestDto
	SiteID *int64 `json:"siteId" validate:"required" error:"required=سایت الزامی است"`
}

// GetByIdProductReviewQuery for retrieving a single product product_review by ID
type GetByIdProductReviewQuery struct {
	ID *int64 `json:"id" validate:"required,gt=0" error:"required=شناسه الزامی است|gt=شناسه باید بزرگتر از 0 باشد"`
}
