package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type SiteRepo struct {
	database *gorm.DB
}

func NewSiteRepository(db *gorm.DB) *SiteRepo {
	return &SiteRepo{
		database: db,
	}
}

func (r *SiteRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Site, int64, error) {
	var sites []domain.Site
	var count int64

	query := r.database.Model(&domain.Site{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&sites)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return sites, count, nil
}

func (r *SiteRepo) GetAllByUserID(userID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Site, int64, error) {
	var sites []domain.Site
	var count int64

	query := r.database.Model(&domain.Site{}).Where("user_id = ?", userID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&sites)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return sites, count, nil
}

func (r *SiteRepo) GetByID(id int64) (domain.Site, error) {
	var site domain.Site
	result := r.database.First(&site, id)
	if result.Error != nil {
		return site, result.Error
	}
	return site, nil
}

func (r *SiteRepo) GetByDomain(domainName string) (domain.Site, error) {
	var site domain.Site
	result := r.database.Where("domain = ?", domainName).First(&site)
	if result.Error != nil {
		return site, result.Error
	}
	return site, nil
}

func (r *SiteRepo) Create(site domain.Site) error {
	result := r.database.Create(&site)
	return result.Error
}

func (r *SiteRepo) Update(site domain.Site) error {
	result := r.database.Save(&site)
	return result.Error
}

func (r *SiteRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Site{}, id)
	return result.Error
}
