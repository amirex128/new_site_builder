package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
	"gorm.io/gorm"
)

type OrderItemRepo struct {
	database *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepo {
	return &OrderItemRepo{
		database: db,
	}
}

func (r *OrderItemRepo) GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.OrderItem], error) {
	var orderItems []domain.OrderItem
	var count int64

	query := r.database.Model(&domain.OrderItem{})
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&orderItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(orderItems, paginationRequestDto, count)
}

func (r *OrderItemRepo) GetAllByOrderID(orderID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.OrderItem], error) {
	var orderItems []domain.OrderItem
	var count int64

	query := r.database.Model(&domain.OrderItem{}).Where("order_id = ?", orderID)
	query.Count(&count)

	limit := paginationRequestDto.PageSize
	offset := (paginationRequestDto.Page - 1) * paginationRequestDto.PageSize

	result := query.Limit(limit).Offset(offset).Find(&orderItems)
	if result.Error != nil {
		return nil, result.Error
	}

	return buildPaginationResponse(orderItems, paginationRequestDto, count)
}

func (r *OrderItemRepo) GetByID(id int64) (*domain.OrderItem, error) {
	var orderItem *domain.OrderItem
	result := r.database.First(&orderItem, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return orderItem, nil
}

func (r *OrderItemRepo) Create(orderItem *domain.OrderItem) error {
	result := r.database.Create(orderItem)
	return result.Error
}

func (r *OrderItemRepo) Update(orderItem *domain.OrderItem) error {
	result := r.database.Save(orderItem)
	return result.Error
}

func (r *OrderItemRepo) Delete(id int64) error {
	result := r.database.Delete(&domain.OrderItem{}, id)
	return result.Error
}

func (r *OrderItemRepo) DeleteByOrderID(orderID int64) error {
	result := r.database.Where("order_id = ?", orderID).Delete(&domain.OrderItem{})
	return result.Error
}
