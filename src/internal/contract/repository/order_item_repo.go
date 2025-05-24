package repository

import (
	common "github.com/amirex128/new_site_builder/src/internal/contract/common"
	"github.com/amirex128/new_site_builder/src/internal/domain"
)

type IOrderItemRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.OrderItem], int64, error)
	GetAllByOrderID(orderID int64, paginationRequestDto common.PaginationRequestDto) (common.PaginationResponseDto[domain.OrderItem], int64, error)
	GetByID(id int64) (domain.OrderItem, error)
	Create(orderItem domain.OrderItem) error
	Update(orderItem domain.OrderItem) error
	Delete(id int64) error
	DeleteByOrderID(orderID int64) error
}
