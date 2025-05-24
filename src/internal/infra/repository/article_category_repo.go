package repository

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

func (r *ArticleCategoryRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.ArticleCategory], error) {
	var categories []domain.ArticleCategory
	var count int64

	query := r.database.Model(&domain.ArticleCategory{}).Where("is_deleted = ?", false)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0
	}

	return buildPaginationResponse(categories, paginationRequestDto, count)
}

func (r *ArticleCategoryRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ArticleCategory, int64, error) {
	var categories []domain.ArticleCategory
	var count int64

	query := r.database.Model(&domain.ArticleCategory{}).
		Where("site_id = ?", siteID).
		Where("is_deleted = ?", false)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return buildPaginationResponse(addresses, paginationRequestDto, count)
}

func (r *ArticleCategoryRepo) GetAllByParentID(parentID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.ArticleCategory, error) {
	var categories []domain.ArticleCategory
	var count int64

	query := r.database.Model(&domain.ArticleCategory{}).
		Where("parent_category_id = ?", parentID).
		Where("is_deleted = ?", false)

	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&categories)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return buildPaginationResponse(addresses, paginationRequestDto, count)
}

func (r *ArticleCategoryRepo) GetByID(id int64) (domain.ArticleCategory, error) {
	var category domain.ArticleCategory

	result := r.database.
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&category)

	if result.Error != nil {
		return domain.ArticleCategory{}, result.Error
	}

	return category, nil
}

func (r *ArticleCategoryRepo) GetBySlug(slug string) (domain.ArticleCategory, error) {
	var category domain.ArticleCategory

	result := r.database.
		Where("slug = ?", slug).
		Where("is_deleted = ?", false).
		First(&category)

	if result.Error != nil {
		return domain.ArticleCategory{}, result.Error
	}

	return category, nil
}

func (r *ArticleCategoryRepo) GetBySlugAndSiteID(slug string, siteID int64) (domain.ArticleCategory, error) {
	var category domain.ArticleCategory

	result := r.database.
		Where("slug = ?", slug).
		Where("site_id = ?", siteID).
		Where("is_deleted = ?", false).
		First(&category)

	if result.Error != nil {
		return domain.ArticleCategory{}, result.Error
	}

	return category, nil
}

func (r *ArticleCategoryRepo) Create(category domain.ArticleCategory) error {
	return r.database.Create(&category).Error
}

func (r *ArticleCategoryRepo) Update(category domain.ArticleCategory) error {
	return r.database.Save(&category).Error
}

func (r *ArticleCategoryRepo) Delete(id int64) error {
	// Soft delete
	return r.database.Model(&domain.ArticleCategory{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

// Media relationship methods

func (r *ArticleCategoryRepo) GetCategoryMedia(categoryID int64) ([]domain.Media, error) {
	var mediaItems []domain.Media

	err := r.database.
		Joins("JOIN category_media ON category_media.media_id = media.id").
		Where("category_media.category_id = ?", categoryID).
		Find(&mediaItems).Error

	if err != nil {
		return nil, err
	}

	return mediaItems, nil
}

func (r *ArticleCategoryRepo) AddMediaToCategory(categoryID int64, mediaID int64) error {
	categoryMedia := domain.ArticleCategoryMedia{
		CategoryID: categoryID,
		MediaID:    mediaID,
	}

	return r.database.Create(&categoryMedia).Error
}

func (r *ArticleCategoryRepo) RemoveMediaFromCategory(categoryID int64, mediaID int64) error {
	return r.database.
		Where("category_id = ? AND media_id = ?", categoryID, mediaID).
		Delete(&domain.ArticleCategoryMedia{}).Error
}

func (r *ArticleCategoryRepo) RemoveAllMediaFromCategory(categoryID int64) error {
	return r.database.
		Where("category_id = ?", categoryID).
		Delete(&domain.ArticleCategoryMedia{}).Error
}
