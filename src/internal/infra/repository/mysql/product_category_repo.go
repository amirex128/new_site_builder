package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type ProductCategoryRepo struct {
	database *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) *ProductCategoryRepo {
	return &ProductCategoryRepo{
		database: db,
	}
}

func (r *ProductCategoryRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.ProductCategory, int64, error) {
	var categories []domain.ProductCategory
	var count int64

	query := r.database.Model(&domain.ProductCategory{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, count, nil
}

func (r *ProductCategoryRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductCategory, int64, error) {
	var categories []domain.ProductCategory
	var count int64

	query := r.database.Model(&domain.ProductCategory{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, count, nil
}

func (r *ProductCategoryRepo) GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ProductCategory, int64, error) {
	var categories []domain.ProductCategory
	var count int64

	query := r.database.Model(&domain.ProductCategory{}).Where("parent_category_id = ?", parentID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, count, nil
}

func (r *ProductCategoryRepo) GetByID(id int64) (domain.ProductCategory, error) {
	var category domain.ProductCategory
	result := r.database.First(&category, id)
	if result.Error != nil {
		return category, result.Error
	}
	return category, nil
}

func (r *ProductCategoryRepo) GetBySlug(slug string) (domain.ProductCategory, error) {
	var category domain.ProductCategory
	result := r.database.Where("slug = ?", slug).First(&category)
	if result.Error != nil {
		return category, result.Error
	}
	return category, nil
}

func (r *ProductCategoryRepo) Create(category domain.ProductCategory) error {
	result := r.database.Create(&category)
	return result.Error
}

func (r *ProductCategoryRepo) Update(category domain.ProductCategory) error {
	result := r.database.Save(&category)
	return result.Error
}

func (r *ProductCategoryRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.ProductCategory{}, id)
	return result.Error
}
