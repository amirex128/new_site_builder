package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type GatewayRepo struct {
	database *gorm.DB
}

func NewGatewayRepository(db *gorm.DB) *GatewayRepo {
	return &GatewayRepo{
		database: db,
	}
}

func (r *GatewayRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Gateway, int64, error) {
	var gateways []domain.Gateway
	var count int64

	query := r.database.Model(&domain.Gateway{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&gateways)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return gateways, count, nil
}

func (r *GatewayRepo) GetBySiteID(siteID int64) (domain.Gateway, error) {
	var gateway domain.Gateway
	result := r.database.Where("site_id = ?", siteID).First(&gateway)
	if result.Error != nil {
		return gateway, result.Error
	}
	return gateway, nil
}

func (r *GatewayRepo) GetByID(id int64) (domain.Gateway, error) {
	var gateway domain.Gateway
	result := r.database.First(&gateway, id)
	if result.Error != nil {
		return gateway, result.Error
	}
	return gateway, nil
}

func (r *GatewayRepo) Create(gateway domain.Gateway) error {
	result := r.database.Create(&gateway)
	return result.Error
}

func (r *GatewayRepo) Update(gateway domain.Gateway) error {
	result := r.database.Save(&gateway)
	return result.Error
}

func (r *GatewayRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Gateway{}, id)
	return result.Error
}
