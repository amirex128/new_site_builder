package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"strings"

	"gorm.io/gorm"
)

type ProvinceRepo struct {
	database *gorm.DB
}

func NewProvinceRepository(db *gorm.DB) *ProvinceRepo {
	return &ProvinceRepo{
		database: db,
	}
}

func (r *ProvinceRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Province], error) {
	var provinces []domain.Province
	var count int64

	query := r.database.Model(&domain.Province{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&provinces)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(provinces, paginationRequestDto, count)
}

func (r *ProvinceRepo) GetByID(id int64) (*domain.Province, error) {
	var province *domain.Province
	result := r.database.First(&province, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return province, nil
}

func (r *ProvinceRepo) GetByName(name string) (*domain.Province, error) {
	var province *domain.Province
	result := r.database.Where("name = ?", name).First(&province)
	if result.Error != nil {
		return nil, result.Error
	}
	return province, nil
}

func (r *ProvinceRepo) Create(province *domain.Province) error {
	result := r.database.Create(province)
	return result.Error
}

func (r *ProvinceRepo) CreateMany(cities []domain.Province) error {
	result := r.database.Create(&cities)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return nil
		}
		return result.Error
	}
	return nil
}

func (r *ProvinceRepo) Update(province *domain.Province) error {
	result := r.database.Save(province)
	return result.Error
}

func (r *ProvinceRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Province{}, id)
	return result.Error
}
