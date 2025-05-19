package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type ProductAttributeRepo struct {
	database *gorm.DB
}

func NewProductAttributeRepository(db *gorm.DB) *ProductAttributeRepo {
	return &ProductAttributeRepo{
		database: db,
	}
}

func (r *ProductAttributeRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.ProductAttribute, int64, error) {
	var attributes []domain.ProductAttribute
	var count int64

	query := r.database.Model(&domain.ProductAttribute{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&attributes)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return attributes, count, nil
}

func (r *ProductAttributeRepo) GetAllByProductID(productID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductAttribute, int64, error) {
	var attributes []domain.ProductAttribute
	var count int64

	query := r.database.Model(&domain.ProductAttribute{}).Where("product_id = ?", productID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&attributes)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return attributes, count, nil
}

func (r *ProductAttributeRepo) GetByID(id int64) (domain.ProductAttribute, error) {
	var attribute domain.ProductAttribute
	result := r.database.First(&attribute, id)
	if result.Error != nil {
		return attribute, result.Error
	}
	return attribute, nil
}

func (r *ProductAttributeRepo) Create(attribute domain.ProductAttribute) error {
	result := r.database.Create(&attribute)
	return result.Error
}

func (r *ProductAttributeRepo) Update(attribute domain.ProductAttribute) error {
	result := r.database.Save(&attribute)
	return result.Error
}

func (r *ProductAttributeRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.ProductAttribute{}, id)
	return result.Error
}
