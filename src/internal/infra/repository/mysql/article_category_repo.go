package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type ArticleCategoryRepo struct {
	database *gorm.DB
}

func NewArticleCategoryRepository(db *gorm.DB) *ArticleCategoryRepo {
	return &ArticleCategoryRepo{
		database: db,
	}
}

func (r *ArticleCategoryRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.BlogCategory, int64, error) {
	var categories []domain.BlogCategory
	var count int64

	query := r.database.Model(&domain.BlogCategory{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, count, nil
}

func (r *ArticleCategoryRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.BlogCategory, int64, error) {
	var categories []domain.BlogCategory
	var count int64

	query := r.database.Model(&domain.BlogCategory{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, count, nil
}

func (r *ArticleCategoryRepo) GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.BlogCategory, int64, error) {
	var categories []domain.BlogCategory
	var count int64

	query := r.database.Model(&domain.BlogCategory{}).Where("parent_category_id = ?", parentID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return categories, count, nil
}

func (r *ArticleCategoryRepo) GetByID(id int64) (domain.BlogCategory, error) {
	var category domain.BlogCategory
	result := r.database.First(&category, id)
	if result.Error != nil {
		return category, result.Error
	}
	return category, nil
}

func (r *ArticleCategoryRepo) GetBySlug(slug string) (domain.BlogCategory, error) {
	var category domain.BlogCategory
	result := r.database.Where("slug = ?", slug).First(&category)
	if result.Error != nil {
		return category, result.Error
	}
	return category, nil
}

func (r *ArticleCategoryRepo) Create(category domain.BlogCategory) error {
	result := r.database.Create(&category)
	return result.Error
}

func (r *ArticleCategoryRepo) Update(category domain.BlogCategory) error {
	result := r.database.Save(&category)
	return result.Error
}

func (r *ArticleCategoryRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.BlogCategory{}, id)
	return result.Error
}
