package repository

import (
	common_contract "go-boilerplate/src/internal/contract/common"
	"go-boilerplate/src/internal/domain"
)

type IProductReviewRepository interface {
	GetAll(paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetAllBySiteID(siteID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetAllByProductID(productID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common_contract.PaginationRequestDto) ([]domain.ProductReview, int64, error)
	GetByID(id int64) (domain.ProductReview, error)
	Create(review domain.ProductReview) error
	Update(review domain.ProductReview) error
	Delete(id int64) error
}
