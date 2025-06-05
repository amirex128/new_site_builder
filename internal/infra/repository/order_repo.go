package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type OrderRepo struct {
	database *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		database: db,
	}
}

func (r *OrderRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Order], error) {
	var orders []domain.Order
	var count int64

	query := r.database.Model(&domain.Order{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(orders, paginationRequestDto, count)
}

func (r *OrderRepo) GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Order], error) {
	var orders []domain.Order
	var count int64

	query := r.database.Model(&domain.Order{}).Where("site_id = ?", siteID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(orders, paginationRequestDto, count)
}

func (r *OrderRepo) GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Order], error) {
	var orders []domain.Order
	var count int64

	query := r.database.Model(&domain.Order{}).Where("customer_id = ?", customerID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(orders, paginationRequestDto, count)
}

func (r *OrderRepo) GetByID(id int64) (*domain.Order, error) {
	var order *domain.Order
	result := r.database.First(&order, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (r *OrderRepo) GetByOrderNumber(orderNumber string) (*domain.Order, error) {
	var order *domain.Order
	result := r.database.Where("order_number = ?", orderNumber).First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (r *OrderRepo) Create(order *domain.Order) error {
	result := r.database.Create(order)
	return result.Error
}

func (r *OrderRepo) Update(order *domain.Order) error {
	result := r.database.Save(order)
	return result.Error
}

func (r *OrderRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.Order{}, id)
	return result.Error
}
