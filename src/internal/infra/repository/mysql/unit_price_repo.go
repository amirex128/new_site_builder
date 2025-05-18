package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type UnitPriceRepo struct {
	database *gorm.DB
}

func NewUnitPriceRepository(db *gorm.DB) *UnitPriceRepo {
	return &UnitPriceRepo{
		database: db,
	}
}

func (r *UnitPriceRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.UnitPrice, int64, error) {
	var unitPrices []domain.UnitPrice
	var count int64

	query := r.database.Model(&domain.UnitPrice{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&unitPrices)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return unitPrices, count, nil
}

func (r *UnitPriceRepo) GetByID(id int64) (domain.UnitPrice, error) {
	var unitPrice domain.UnitPrice
	result := r.database.First(&unitPrice, id)
	if result.Error != nil {
		return unitPrice, result.Error
	}
	return unitPrice, nil
}

func (r *UnitPriceRepo) GetByName(name string) (domain.UnitPrice, error) {
	var unitPrice domain.UnitPrice
	result := r.database.Where("name = ?", name).First(&unitPrice)
	if result.Error != nil {
		return unitPrice, result.Error
	}
	return unitPrice, nil
}

func (r *UnitPriceRepo) Create(unitPrice domain.UnitPrice) error {
	result := r.database.Create(&unitPrice)
	return result.Error
}

func (r *UnitPriceRepo) Update(unitPrice domain.UnitPrice) error {
	result := r.database.Save(&unitPrice)
	return result.Error
}

func (r *UnitPriceRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.UnitPrice{}, id)
	return result.Error
}
