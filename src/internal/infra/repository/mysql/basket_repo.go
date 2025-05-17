package mysql

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"

	"gorm.io/gorm"
)

type BasketRepo struct {
	database *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepo {
	return &BasketRepo{
		database: db,
	}
}

func (r *BasketRepo) GetAll(paginationRequestDto common.PaginationRequestDto) ([]domain.Basket, int64, error) {
	var baskets []domain.Basket
	var count int64

	query := r.database.Model(&domain.Basket{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&baskets)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return baskets, count, nil
}

func (r *BasketRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Basket, int64, error) {
	var baskets []domain.Basket
	var count int64

	query := r.database.Model(&domain.Basket{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&baskets)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return baskets, count, nil
}

func (r *BasketRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) ([]domain.Basket, int64, error) {
	var baskets []domain.Basket
	var count int64

	query := r.database.Model(&domain.Basket{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&baskets)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return baskets, count, nil
}

func (r *BasketRepo) GetByID(id int64) (domain.Basket, error) {
	var basket domain.Basket
	result := r.database.First(&basket, id)
	if result.Error != nil {
		return basket, result.Error
	}
	return basket, nil
}

func (r *BasketRepo) GetActiveBasketByCustomerID(customerID int64) (domain.Basket, error) {
	var basket domain.Basket
	result := r.database.Where("customer_id = ? AND status = 'active'", customerID).First(&basket)
	if result.Error != nil {
		return basket, result.Error
	}
	return basket, nil
}

func (r *BasketRepo) Create(basket domain.Basket) error {
	result := r.database.Create(&basket)
	return result.Error
}

func (r *BasketRepo) Update(basket domain.Basket) error {
	result := r.database.Save(&basket)
	return result.Error
}

func (r *BasketRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Basket{}, id)
	return result.Error
}
