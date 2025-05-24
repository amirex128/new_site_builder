package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type HeaderFooterRepo struct {
	database *gorm.DB
}

func NewHeaderFooterRepository(db *gorm.DB) *HeaderFooterRepo {
	return &HeaderFooterRepo{
		database: db,
	}
}

func (r *HeaderFooterRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.HeaderFooter], error) {
	var headerFooters []domain.HeaderFooter
	var count int64

	query := r.database.Model(&domain.HeaderFooter{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&headerFooters)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(headerFooters, paginationRequestDto, count)
}

func (r *HeaderFooterRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.HeaderFooter], error) {
	var headerFooters []domain.HeaderFooter
	var count int64

	query := r.database.Model(&domain.HeaderFooter{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&headerFooters)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(headerFooters, paginationRequestDto, count)
}

func (r *HeaderFooterRepo) GetByID(id int64) (domain.HeaderFooter, error) {
	var headerFooter domain.HeaderFooter
	result := r.database.First(&headerFooter, id)
	if result.Error != nil {
		return headerFooter, result.Error
	}
	return headerFooter, nil
}

func (r *HeaderFooterRepo) Create(headerFooter domain.HeaderFooter) error {
	result := r.database.Create(&headerFooter)
	return result.Error
}

func (r *HeaderFooterRepo) Update(headerFooter domain.HeaderFooter) error {
	result := r.database.Save(&headerFooter)
	return result.Error
}

func (r *HeaderFooterRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.HeaderFooter{}, id)
	return result.Error
}
