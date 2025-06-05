package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type ProductReviewRepo struct {
	database *gorm.DB
}

func NewProductReviewRepository(db *gorm.DB) *ProductReviewRepo {
	return &ProductReviewRepo{
		database: db,
	}
}

func (r *ProductReviewRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductReview], error) {
	var reviews []domain.ProductReview
	var count int64

	query := r.database.Model(&domain.ProductReview{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(reviews, paginationRequestDto, count)
}

func (r *ProductReviewRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductReview], error) {
	var reviews []domain.ProductReview
	var count int64

	query := r.database.Model(&domain.ProductReview{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(reviews, paginationRequestDto, count)
}

func (r *ProductReviewRepo) GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductReview], error) {
	var reviews []domain.ProductReview
	var count int64

	query := r.database.Model(&domain.ProductReview{}).Where("product_id = ?", productID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(reviews, paginationRequestDto, count)
}

func (r *ProductReviewRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.ProductReview], error) {
	var reviews []domain.ProductReview
	var count int64

	query := r.database.Model(&domain.ProductReview{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(reviews, paginationRequestDto, count)
}

func (r *ProductReviewRepo) GetByID(id int64) (*domain.ProductReview, error) {
	var review domain.ProductReview
	result := r.database.First(&review, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (r *ProductReviewRepo) Create(review *domain.ProductReview) error {
	result := r.database.Create(review)
	return result.Error
}

func (r *ProductReviewRepo) Update(review *domain.ProductReview) error {
	result := r.database.Save(review)
	return result.Error
}

func (r *ProductReviewRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.ProductReview{}, id)
	return result.Error
}
