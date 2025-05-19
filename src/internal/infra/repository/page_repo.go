package repository

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

func (r *PageRepo) AddMediaToPage(pageID int64, mediaIDs []int64) error {
	// Begin a transaction
	tx := r.database.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete existing media associations for this page
	if err := tx.Exec("DELETE FROM page_media WHERE page_id = ?", pageID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new media associations
	for _, mediaID := range mediaIDs {
		if err := tx.Exec("INSERT INTO page_media (page_id, media_id) VALUES (?, ?)", pageID, mediaID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *PageRepo) RemoveAllMediaFromPage(pageID int64) error {
	return r.database.Exec("DELETE FROM page_media WHERE page_id = ?", pageID).Error
}

func (r *PageRepo) GetByIDAndSiteID(id, siteID int64) (domain.Page, error) {
	var page domain.Page
	result := r.database.Where("id = ? AND site_id = ?", id, siteID).First(&page)
	if result.Error != nil {
		return page, result.Error
	}
	return page, nil
}

func (r *PageRepo) GetByIDs(ids []int64, siteID int64) ([]domain.Page, error) {
	var pages []domain.Page
	result := r.database.Where("id IN ? AND site_id = ?", ids, siteID).Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pages, nil
}

func (r *PageRepo) GetByPaths(paths []string, siteID int64) ([]domain.Page, error) {
	var pages []domain.Page
	result := r.database.Where("slug IN ? AND site_id = ?", paths, siteID).Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pages, nil
}
