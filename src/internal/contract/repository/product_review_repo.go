package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IProductReviewRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetByID(id int64) (domain.ProductReview, error)
	Create(review domain.ProductReview) error
	Update(review domain.ProductReview) error
	Delete(id int64) error
}
