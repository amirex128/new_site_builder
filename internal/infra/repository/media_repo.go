package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type MediaRepo struct {
	database *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepo {
	return &MediaRepo{
		database: db,
	}
}

func (r *MediaRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Media], error) {
	var media []domain.Media
	var count int64

	query := r.database.Model(&domain.Media{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&media)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(media, paginationRequestDto, count)
}

func (r *MediaRepo) GetByID(id int64) (*domain.Media, error) {
	var media *domain.Media
	result := r.database.First(&media, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return media, nil
}

func (r *MediaRepo) Create(media *domain.Media) error {
	result := r.database.Create(media)
	return result.Error
}

func (r *MediaRepo) Update(media *domain.Media) error {
	result := r.database.Save(media)
	return result.Error
}

func (r *MediaRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Media{}, id)
	return result.Error
}
