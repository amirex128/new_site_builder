package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type CityRepo struct {
	database *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepo {
	return &CityRepo{
		database: db,
	}
}

func (r *CityRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.City, int64, error) {
	var cities []domain.City
	var count int64

	query := r.database.Model(&domain.City{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&cities)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return cities, count, nil
}

func (r *CityRepo) GetByID(id int64) (domain.City, error) {
	var city domain.City
	result := r.database.First(&city, id)
	if result.Error != nil {
		return city, result.Error
	}
	return city, nil
}

func (r *CityRepo) GetByName(name string) (domain.City, error) {
	var city domain.City
	result := r.database.Where("name = ?", name).First(&city)
	if result.Error != nil {
		return city, result.Error
	}
	return city, nil
}

func (r *CityRepo) Create(city domain.City) error {
	result := r.database.Create(&city)
	return result.Error
}

func (r *CityRepo) Update(city domain.City) error {
	result := r.database.Save(&city)
	return result.Error
}

func (r *CityRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.City{}, id)
	return result.Error
}
