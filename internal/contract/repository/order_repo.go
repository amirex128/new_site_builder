package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
	"github.com/amirex128/new_site_builder/internal/domain"
)

type IOrderRepository interface {
	GetAll(paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Order], error)
	GetAllBySiteID(siteID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Order], error)
	GetAllByCustomerID(customerID int64, paginationRequestDto common.PaginationRequestDto) (*common.PaginationResponseDto[domain.Order], error)
	GetByID(id int64) (*domain.Order, error)
	GetByOrderNumber(orderNumber string) (*domain.Order, error)
	Create(order *domain.Order) error
	Update(order *domain.Order) error
	Delete(id int64) error
}
