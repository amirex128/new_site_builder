package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

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

func (r *ProvinceRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Province, int64, error) {
	var provinces []domain.Province
	var count int64

	query := r.database.Model(&domain.Province{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&provinces)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return provinces, count, nil
}

func (r *ProvinceRepo) GetByID(id int64) (domain.Province, error) {
	var province domain.Province
	result := r.database.First(&province, id)
	if result.Error != nil {
		return province, result.Error
	}
	return province, nil
}

func (r *ProvinceRepo) GetByName(name string) (domain.Province, error) {
	var province domain.Province
	result := r.database.Where("name = ?", name).First(&province)
	if result.Error != nil {
		return province, result.Error
	}
	return province, nil
}

func (r *ProvinceRepo) Create(province domain.Province) error {
	result := r.database.Create(&province)
	return result.Error
}

func (r *ProvinceRepo) Update(province domain.Province) error {
	result := r.database.Save(&province)
	return result.Error
}

func (r *ProvinceRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Province{}, id)
	return result.Error
}
