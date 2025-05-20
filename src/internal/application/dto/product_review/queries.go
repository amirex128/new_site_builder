package product_review

import (
	"github.com/amirex128/new_site_builder/src/internal/contract/common"
)

// AdminGetAllProductReviewQuery for admin listing of all article reviews with pagination
type AdminGetAllProductReviewQuery struct {
	common.PaginationRequestDto
}

// GetAllProductReviewQuery for listing article reviews by site ID with pagination
type GetAllProductReviewQuery struct {
	common.PaginationRequestDto
	ProductID *int64 `json:"productId" form:"productId" validate:"required,gt=0"`
	SiteID    *int64 `json:"siteId" form:"siteId" validate:"required,gt=0"`
}

// GetByIdProductReviewQuery for retrieving a single article product_review by ID
type GetByIdProductReviewQuery struct {
	ID *int64 `json:"id" form:"id" validate:"required,gt=0"`
}
