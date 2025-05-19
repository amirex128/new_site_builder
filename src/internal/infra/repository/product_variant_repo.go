package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"time"

	"gorm.io/gorm"
)

type ProductVariantRepo struct {
	database *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) *ProductVariantRepo {
	return &ProductVariantRepo{
		database: db,
	}
}

func (r *ProductVariantRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.ProductVariant, int64, error) {
	var variants []domain.ProductVariant
	var count int64

	query := r.database.Model(&domain.ProductVariant{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&variants)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return variants, count, nil
}

func (r *ProductVariantRepo) GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductVariant, int64, error) {
	var variants []domain.ProductVariant
	var count int64

	query := r.database.Model(&domain.ProductVariant{}).Where("product_id = ?", productID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&variants)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return variants, count, nil
}

func (r *ProductVariantRepo) GetByID(id int64) (domain.ProductVariant, error) {
	var variant domain.ProductVariant
	result := r.database.First(&variant, id)
	if result.Error != nil {
		return variant, result.Error
	}
	return variant, nil
}

func (r *ProductVariantRepo) Create(variant domain.ProductVariant) error {
	result := r.database.Create(&variant)
	return result.Error
}

func (r *ProductVariantRepo) Update(variant domain.ProductVariant) error {
	result := r.database.Save(&variant)
	return result.Error
}

func (r *ProductVariantRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.ProductVariant{}, id)
	return result.Error
}

func (r *ProductVariantRepo) DecreaseStock(variantID int64, quantity int) error {
	// Get the current variant
	var variant domain.ProductVariant
	if err := r.database.First(&variant, variantID).Error; err != nil {
		return err
	}

	// Check if there's enough stock
	if variant.Stock < quantity {
		return gorm.ErrInvalidData
	}

	// Decrease the stock
	variant.Stock -= quantity
	variant.UpdatedAt = time.Now()

	// Update the variant
	return r.database.Save(&variant).Error
}
