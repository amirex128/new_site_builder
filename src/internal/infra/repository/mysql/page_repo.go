package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type PageRepo struct {
	database *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepo {
	return &PageRepo{
		database: db,
	}
}

func (r *PageRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Page, int64, error) {
	var pages []domain.Page
	var count int64

	query := r.database.Model(&domain.Page{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&pages)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return pages, count, nil
}

func (r *PageRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Page, int64, error) {
	var pages []domain.Page
	var count int64

	query := r.database.Model(&domain.Page{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&pages)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return pages, count, nil
}

func (r *PageRepo) GetByID(id int64) (domain.Page, error) {
	var page domain.Page
	result := r.database.First(&page, id)
	if result.Error != nil {
		return page, result.Error
	}
	return page, nil
}

func (r *PageRepo) GetBySlug(slug string, siteID int64) (domain.Page, error) {
	var page domain.Page
	result := r.database.Where("slug = ? AND site_id = ?", slug, siteID).First(&page)
	if result.Error != nil {
		return page, result.Error
	}
	return page, nil
}

func (r *PageRepo) Create(page domain.Page) error {
	result := r.database.Create(&page)
	return result.Error
}

func (r *PageRepo) Update(page domain.Page) error {
	result := r.database.Save(&page)
	return result.Error
}

func (r *PageRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Page{}, id)
	return result.Error
}
